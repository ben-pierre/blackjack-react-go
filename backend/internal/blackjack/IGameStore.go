package blackjack

import (
	"blackjack/internal/domain"
	"uuid"
)

type GameStore interface {
	GetGame(gameId uuid.UUID, userId uuid.UUID, guestToken string) (domain.Game, error)
	GetGames(userId uuid.UUID, limit, offset int) ([]domain.Game, error)
	UpdateGame(game domain.Game) error
}
