package blackjack

import (
	"blackjack/internal/domain"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"uuid"
)

type Service struct {
	db GameStore
}

func New(pg GameStore) *Service {
	return &Service{
		db: pg,
	}
}

func (s *Service) HandleGameAction(gameUUID uuid.UUID, userUUID uuid.UUID, guestToken string, action func(*domain.Game) error) (domain.Game, error) {
	game, err := s.db.GetGame(gameUUID, userUUID, guestToken)
	if err != nil {
		return game, fmt.Errorf("game not found : %v", err)
	}
	err = action(&game)
	if err != nil {
		return domain.Game{}, err
	}
	err = s.db.UpdateGame(game)
	return game, err
}

func (s *Service) GetNewGame(userID uuid.UUID) (domain.Game, error) {
	newGameUUID := uuid.New()
	newGame := domain.Game{GameId: newGameUUID, Deck: domain.GetNewDeck(domain.DeckConfig{Shuffle: true}), PlayerHand: domain.Hand{}, DealerHand: domain.Hand{}, Phase: "player"}

	if userID != uuid.Nil() {
		newGame.UserId = userID
	} else {
		token, err := generateGuestToken(24)
		if err != nil {
			return domain.Game{}, err
		}
		newGame.GuestToken = token
	}

	newGame.PlayerHand.Cards = append(newGame.PlayerHand.Cards, newGame.Deck.DrawCard())
	newGame.DealerHand.Cards = append(newGame.DealerHand.Cards, newGame.Deck.DrawCard())
	newGame.PlayerHand.Cards = append(newGame.PlayerHand.Cards, newGame.Deck.DrawCard())
	newGame.DealerHand.Cards = append(newGame.DealerHand.Cards, newGame.Deck.DrawCard())
	newGame.CheckWin()
	err := s.db.UpdateGame(newGame)
	return newGame, err
}

func generateGuestToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
