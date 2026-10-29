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
	"runtime/debug"
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

	userId, _ := r.Context().Value(userIDKey{}).(uuid.UUID) //checked in auth middle ware

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

func getGame(w http.ResponseWriter, r *http.Request, service *blackjack.Service) {
	userId, ok := r.Context().Value(userIDKey{}).(uuid.UUID)
	if !ok {
		http.Error(w, "Unathorized", http.StatusUnauthorized)
	}
	gameId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid Game UUID", http.StatusBadRequest)
		return
	}
	game, err := service.GetGame(gameId, userId)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(game.ToResponse())
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
		http.Error(w, "Invalid Game UUID", http.StatusBadRequest)
		return
	}

	game, gameErr := service.HandleGameAction(gameId, userID, guestToken, action)

	if gameErr != nil {
		log.Printf("Failed to handle game action: %v", gameErr.Error())
		http.Error(w, gameErr.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(game.ToResponse())
}

func getUserAccount(w http.ResponseWriter, r *http.Request, service *blackjack.Service) {
	w.Header().Set("Content-Type", "application/json")
	userId, ok := r.Context().Value(userIDKey{}).(uuid.UUID)
	if !ok || userId == uuid.Nil() {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := service.GetUser(userId)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	games, err := service.GetGames(user.UserId, 10, 0)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	userAccount := domain.UserAccountResponse{Balance: user.Balance, Games: games}
	json.NewEncoder(w).Encode(userAccount)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	log.Println("Initializing PostgreSQL service")
	pgService := postgres.New()
	log.Println("PostgreSQL service initialized")
	log.Println("Initializing Blackjack service")
	service := blackjack.New(pgService, pgService)
	log.Println("Blackjack service initialized")
	defer pgService.Close()

	jwtIssuer := os.Getenv("COGNITOISSUER")
	jwksURL := jwtIssuer + "/.well-known/jwks.json"
	log.Printf("Initializing Cognito JWKS client: %s", jwksURL)
	keyFunc, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		log.Fatalf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
	}
	log.Println("Cognito JWKS client initialized")

	mux := http.NewServeMux()
	log.Println("Registering endpoints")
	mux.Handle("POST /blackjack/new",
		chainMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				newGame(w, r, service)
			}),
			func(next http.Handler) http.Handler {
				return authMiddleware(next, keyFunc, true)
			},
			logMiddleware,
		),
	)

	mux.Handle("POST /blackjack/{id}/hit",
		chainMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handleGameAction(w, r, (*domain.Game).Hit, service)
			}),
			func(next http.Handler) http.Handler {
				return authMiddleware(next, keyFunc, true)
			},
			logMiddleware,
		),
	)

	mux.Handle("POST /blackjack/{id}/stay",
		chainMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handleGameAction(w, r, (*domain.Game).Stay, service)
			}),
			func(next http.Handler) http.Handler {
				return authMiddleware(next, keyFunc, true)
			},
			logMiddleware,
		),
	)

	mux.Handle("GET /blackjack/{id}",
		chainMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				getGame(w, r, service)
			}),
			func(next http.Handler) http.Handler {
				return authMiddleware(next, keyFunc, true)
			},
			logMiddleware,
		),
	)

	mux.Handle("GET /user/account",
		chainMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				getUserAccount(w, r, service)
			}),
			func(next http.Handler) http.Handler {
				return authMiddleware(next, keyFunc, true)
			},
			logMiddleware,
		),
	)

	mux.HandleFunc("GET /healthz", healthcheck)

	log.Println("End points registered\nBuilding handler")

	handler := chainMiddleware(
		mux,
		recoverMiddleware,
		corsMiddleware,
	)

	log.Println("Built handler\nStarting app")

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Lambda handler starting")
		lambda.Start(httpadapter.NewV2(handler).ProxyWithContext)
	} else {
		log.Print("Started server")
		log.Fatal(http.ListenAndServe(":8080", handler))
	}
}

func chainMiddleware(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Guest-Token")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)
		// userID, _ := r.Context().Value(userIDKey{}).(uuid.UUID)

		// guestToken := r.Header.Get("X-Guest-Token")
		// log.Printf("Handling request with UUID : %v, guestToken : %v", userID, guestToken)
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
