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
	INSERT INTO games (gameId, phase, wager, deck, playerhand, dealerhand, userId, guestToken)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (gameId) DO UPDATE SET
		phase = EXCLUDED.phase,
		wager = EXCLUDED.wager,
		deck = EXCLUDED.deck,
		playerhand = EXCLUDED.playerhand,
		dealerhand = EXCLUDED.dealerhand,
		userId = EXCLUDED.userId,
		guestToken = EXCLUDED.guestToken;`

	playerHand, err := json.Marshal(game.PlayerHand)
	if err != nil {
		return fmt.Errorf("failed to marshal PlayerHand %s: %w", game.GameId, err)
	}

	dealerHand, err := json.Marshal(game.DealerHand)
	if err != nil {
		return fmt.Errorf("failed to marshal DealerHand %s: %w", game.GameId, err)
	}

	deck, err := json.Marshal(game.Deck)
	if err != nil {
		return fmt.Errorf("failed to marshal deck %s: %w", game.GameId, err)
	}

	_, err = s.db.Exec(sqlStatement, game.GameId, game.Phase, game.Wager, deck, playerHand, dealerHand, game.UserId, game.GuestToken)
	if err != nil {
		return fmt.Errorf("failed to store game %s: %w", game.GameId, err)
	}

	return nil
}

func (s *Service) GetGame(gameId uuid.UUID, userId uuid.UUID, guestToken string) (domain.Game, error) {
	var game domain.Game
	var deck, playerHand, dealerHand []byte

	err := s.db.QueryRow(`
		SELECT gameId, phase, wager, deck, playerhand, dealerhand, userId, guestToken
		FROM games
		WHERE gameId = $1
		AND (userId = $2 or guestToken=$3)
	`, gameId, userId, guestToken).Scan(
		&game.GameId,
		&game.Phase,
		&game.Wager,
		&deck,
		&playerHand,
		&dealerHand,
		&game.UserId,
		&game.GuestToken,
	)
	if err != nil {
		return domain.Game{}, fmt.Errorf("failed to get game %s: %w", gameId, err)
	}

	if err := json.Unmarshal(deck, &game.Deck); err != nil {
		return domain.Game{}, fmt.Errorf("failed to unmarshal deck for game %s: %w", gameId, err)
	}

	if err := json.Unmarshal(playerHand, &game.PlayerHand); err != nil {
		return domain.Game{}, fmt.Errorf("failed to unmarshal player hand for game %s: %w", gameId, err)
	}

	if err := json.Unmarshal(dealerHand, &game.DealerHand); err != nil {
		return domain.Game{}, fmt.Errorf("failed to unmarshal dealer hand for game %s: %w", gameId, err)
	}

	return game, nil
}

func (s *Service) GetGames(userId uuid.UUID, limit, offset int) ([]domain.Game, error) {
	rows, err := s.db.Query(`
		SELECT gameId, phase, wager, deck, playerhand, dealerhand, userId, guestToken
		FROM games
		WHERE userId = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get games for user %s: %w", userId, err)
	}
	defer rows.Close()

	var games []domain.Game

	for rows.Next() {
		var game domain.Game
		var deck, playerHand, dealerHand []byte

		if err := rows.Scan(
			&game.GameId,
			&game.Phase,
			&game.Wager,
			&deck,
			&playerHand,
			&dealerHand,
			&game.UserId,
			&game.GuestToken,
		); err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}

		if err := json.Unmarshal(deck, &game.Deck); err != nil {
			return nil, fmt.Errorf("failed to unmarshal deck for game %s: %w", game.GameId, err)
		}
		if err := json.Unmarshal(playerHand, &game.PlayerHand); err != nil {
			return nil, fmt.Errorf("failed to unmarshal player hand for game %s: %w", game.GameId, err)
		}
		if err := json.Unmarshal(dealerHand, &game.DealerHand); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dealer hand for game %s: %w", game.GameId, err)
		}

		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while iterating games: %w", err)
	}

	return games, nil
}

func (s *Service) CreateUser(user domain.User) error {
	_, err := s.db.Exec(`
        INSERT INTO users (user_id, balance)
        VALUES ($1, $2)
        ON CONFLICT (user_id) DO NOTHING
    `, user.UserId, user.Balance)

	return err
}

func (s *Service) UpdateUser(user domain.User) error {
	_, err := s.db.Exec(`
        UPDATE users
        SET balance = $2
        WHERE user_id = $1
    `, user.UserId, user.Balance)

	return err
}

func (s *Service) GetUser(userId uuid.UUID) (*domain.User, error) {
	var user domain.User

	err := s.db.QueryRow(`
        SELECT user_id, balance, created_at
        FROM users
        WHERE user_id = $1
    `, userId).Scan(&user.UserId, &user.Balance, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
