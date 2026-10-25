package blackjack

import (
	"blackjack/internal/domain"
	"fmt"
	"sync"
	"uuid"
)

type Service struct {
	mu sync.RWMutex
	db GameStore
}

func New(pg GameStore) *Service {
	return &Service{
		db: pg,
	}
}

func (s *Service) HandleGameAction(uuid uuid.UUID, action func(*domain.Game) error) (domain.Game, error) {
	s.mu.RLock()
	game, err := s.db.GetGame(uuid)
	s.mu.RUnlock()
	if err != nil {
		return game, fmt.Errorf("game not found : %v", err)
	}
	gameErr := action(&game)
	s.mu.Lock()
	s.db.UpdateGame(game)
	s.mu.Unlock()
	return game, gameErr
}

func (s *Service) GetNewGame() domain.Game {
	newGameUUID := uuid.New()
	newGame := domain.Game{UUID: newGameUUID, Deck: domain.GetNewDeck(domain.DeckConfig{Shuffle: true}), PlayerHand: domain.Hand{}, DealerHand: domain.Hand{}, Phase: "player"}
	newGame.PlayerHand.Cards = append(newGame.PlayerHand.Cards, newGame.Deck.DrawCard())
	newGame.DealerHand.Cards = append(newGame.DealerHand.Cards, newGame.Deck.DrawCard())
	newGame.PlayerHand.Cards = append(newGame.PlayerHand.Cards, newGame.Deck.DrawCard())
	newGame.DealerHand.Cards = append(newGame.DealerHand.Cards, newGame.Deck.DrawCard())
	newGame.CheckWin()
	s.mu.Lock()
	s.db.UpdateGame(newGame)
	s.mu.Unlock()
	return newGame
}
