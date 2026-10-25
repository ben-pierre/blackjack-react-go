package postgres

import (
	"blackjack/internal/domain"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"uuid"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("PGHOST")
	port     = 5432
	user     = os.Getenv("PGUSER")
	password = os.Getenv("PGPASSWORD")
	dbname   = os.Getenv("PGDATABASE")
)

type Service struct {
	healthy bool
	db      *sql.DB
}

func New() *Service {
	healthy := false
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully connected to %v!\n", host)
	healthy = true
	return &Service{
		healthy: healthy,
		db:      db,
	}
}

func (s *Service) Close() error {
	return s.db.Close()
}

func (s *Service) UpdateGame(game domain.Game) error {
	sqlStatement := `
	INSERT INTO games (uuid, phase, wager, deck, playerhand, dealerhand)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (uuid) DO UPDATE SET
		phase = EXCLUDED.phase,
		wager = EXCLUDED.wager,
		deck = EXCLUDED.deck,
		playerhand = EXCLUDED.playerhand,
		dealerhand = EXCLUDED.dealerhand;`

	playerHand, err := json.Marshal(game.PlayerHand)
	if err != nil {
		return fmt.Errorf("failed to marshall PlayerHand %s: %w", game.UUID, err)
	}

	dealerHand, err := json.Marshal(game.DealerHand)
	if err != nil {
		return fmt.Errorf("failed to marshall DealerHand %s: %w", game.UUID, err)
	}

	deck, err := json.Marshal(game.Deck)
	if err != nil {
		return fmt.Errorf("failed to marshall deck %s: %w", game.UUID, err)
	}

	_, err = s.db.Exec(sqlStatement, game.UUID, game.Phase, game.Wager, deck, playerHand, dealerHand)
	if err != nil {
		return fmt.Errorf("failed to store game %s: %w", game.UUID, err)
	}

	return nil
}

func (s *Service) GetGame(id uuid.UUID) (domain.Game, error) {
	var game domain.Game
	var deck, playerHand, dealerHand []byte

	err := s.db.QueryRow(`
		SELECT uuid, phase, wager, deck, playerhand, dealerhand
		FROM games
		WHERE uuid = $1
	`, id).Scan(
		&game.UUID,
		&game.Phase,
		&game.Wager,
		&deck,
		&playerHand,
		&dealerHand,
	)
	if err != nil {
		return domain.Game{}, fmt.Errorf("failed to get game %s: %w", id, err)
	}

	if err := json.Unmarshal(deck, &game.Deck); err != nil {
		return domain.Game{}, fmt.Errorf("failed to unmarshal deck for game %s: %w", id, err)
	}

	if err := json.Unmarshal(playerHand, &game.PlayerHand); err != nil {
		return domain.Game{}, fmt.Errorf("failed to unmarshal player hand for game %s: %w", id, err)
	}

	if err := json.Unmarshal(dealerHand, &game.DealerHand); err != nil {
		return domain.Game{}, fmt.Errorf("failed to unmarshal dealer hand for game %s: %w", id, err)
	}

	return game, nil
}
