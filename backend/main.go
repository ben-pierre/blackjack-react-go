package main

import (
	"blackjack/internal/blackjack"
	"blackjack/internal/domain"
	"blackjack/internal/postgres"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"uuid"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func newGame(w http.ResponseWriter, r *http.Request, service *blackjack.Service) {
	newGame := service.GetNewGame()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGame.ToResponse())
}

// func getGames(w http.ResponseWriter, r *http.Request, service *blackjack.Service) {
// 	summaries := service.GetGames()
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(summaries)

// }

func handleGameAction(
	w http.ResponseWriter,
	r *http.Request,
	action func(*domain.Game) error,
	service *blackjack.Service,
) {
	w.Header().Set("Content-Type", "application/json")

	gameId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Invalid UUID",
		})
		return
	}

	game, gameErr := service.HandleGameAction(gameId, action)

	if gameErr != nil {
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
	mux := http.NewServeMux()

	mux.HandleFunc("POST /blackjack/new", func(w http.ResponseWriter, r *http.Request) {
		newGame(w, r, service)
	})
	// mux.HandleFunc("GET /blackjack/games", func(w http.ResponseWriter, r *http.Request) { getGames(w, r, service) })

	mux.HandleFunc("POST /blackjack/{id}/hit", func(w http.ResponseWriter, r *http.Request) {
		handleGameAction(w, r, (*domain.Game).Hit, service)
	})

	mux.HandleFunc("POST /blackjack/{id}/stay", func(w http.ResponseWriter, r *http.Request) {
		handleGameAction(w, r, (*domain.Game).Stay, service)
	})

	mux.HandleFunc("GET /healthz", healthcheck)

	handler := corsMiddleware(mux)
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Lambda handler starting")
		lambda.Start(httpadapter.NewV2(handler).ProxyWithContext)
	} else {
		log.Print("Started server")
		log.Fatal(http.ListenAndServe(":8080", handler))
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
