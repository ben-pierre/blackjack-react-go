package main

import (
	"blackjack/internal/blackjack"
	"blackjack/internal/domain"
	"blackjack/internal/postgres"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"uuid"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/golang-jwt/jwt/v5"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type userIDKey struct{}

func newGame(w http.ResponseWriter, r *http.Request, service *blackjack.Service) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value(userIDKey{}).(uuid.UUID)
	if !ok {
		userId = uuid.Nil()
		log.Printf("Creating guest game")
	} else {
		log.Printf("Creating new game for user %v", userId)
	}

	newGame, err := service.GetNewGame(userId)

	if err != nil {
		log.Printf("Failed to create new game: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Failed to create game",
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGame.ToResponse())
}

func handleGameAction(
	w http.ResponseWriter,
	r *http.Request,
	action func(*domain.Game) error,
	service *blackjack.Service,
) {
	userID, ok := r.Context().Value(userIDKey{}).(uuid.UUID)

	guestToken := r.Header.Get("X-Guest-Token")
	if !ok {
		userID = uuid.Nil()
	}

	log.Printf("Processing game action for user %v", userID)

	w.Header().Set("Content-Type", "application/json")
	gameId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Invalid Game UUID",
		})
		return
	}

	game, gameErr := service.HandleGameAction(gameId, userID, guestToken, action)

	if gameErr != nil {
		log.Printf("Failed to handle game action: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: gameErr.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(game.ToResponse())
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	pgService := postgres.New()
	var service = blackjack.New(pgService)
	defer pgService.Close()

	// Create the keyfunc.Keyfunc.
	jwtIssuer := os.Getenv("COGNITOISSUER")
	jwksURL := jwtIssuer + "/.well-known/jwks.json"
	keyFunc, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL}) // Context is used to end the refresh goroutine.
	if err != nil {
		log.Fatalf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
	}

	mux := http.NewServeMux()

	mux.Handle("POST /blackjack/new", authMiddleware(logMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newGame(w, r, service)
		})), keyFunc, true,
	))

	mux.Handle("POST /blackjack/{id}/hit",
		authMiddleware(
			logMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handleGameAction(w, r, (*domain.Game).Hit, service)
				})), keyFunc, true,
		))

	mux.Handle("POST /blackjack/{id}/stay",
		authMiddleware(
			logMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handleGameAction(w, r, (*domain.Game).Stay, service)
				}),
			),
			keyFunc,
			true,
		),
	)

	mux.HandleFunc("GET /healthz", healthcheck)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Lambda handler starting")
		lambda.Start(httpadapter.NewV2(mux).ProxyWithContext)
	} else {
		log.Print("Started server")
		log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, _ := r.Context().Value(userIDKey{}).(uuid.UUID)

		guestToken := r.Header.Get("X-Guest-Token")
		log.Printf("Handling request with UUID : %v, guestToken : %v", userID, guestToken)
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(
	next http.Handler,
	keyFunc keyfunc.Keyfunc,
	optional bool,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" && optional {
			next.ServeHTTP(w, r)
			log.Print("Optional auth set, no bearer token passed, skipping")
			return
		}
		log.Printf("Found auth header %v, passing", authHeader)

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		if bearerToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parsed, err := jwt.Parse(
			bearerToken,
			keyFunc.Keyfunc,
			jwt.WithValidMethods([]string{"RS256"}),
			jwt.WithIssuer(os.Getenv("COGNITOISSUER")),
		)
		if err != nil || parsed == nil || !parsed.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenUse, _ := claims["token_use"].(string)
		clientID, _ := claims["client_id"].(string)

		if tokenUse != "access" || clientID != os.Getenv("COGNITOCLIENTID") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := parsed.Claims.GetSubject()
		if err != nil || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userUUID, err := uuid.Parse(userID)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey{}, userUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
