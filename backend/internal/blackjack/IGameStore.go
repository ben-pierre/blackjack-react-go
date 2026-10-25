package blackjack

import (
	"blackjack/internal/domain"
	"uuid"
)

type GameStore interface {
	GetGame(id uuid.UUID) (domain.Game, error)
	UpdateGame(game domain.Game) error
}
