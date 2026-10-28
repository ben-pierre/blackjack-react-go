package blackjack

import (
	"blackjack/internal/domain"
	"uuid"
)

type UserStore interface {
	CreateUser(user domain.User) error
	GetUser(userId uuid.UUID) (*domain.User, error)
	UpdateUser(user domain.User) error
}
