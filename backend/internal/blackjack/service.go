package blackjack

import (
	"blackjack/internal/domain"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"uuid"
)

var defaultUserBalance int = 50

type Service struct {
	gs GameStore
	us UserStore
}

func New(gs GameStore, us UserStore) *Service {
	return &Service{
		gs: gs,
		us: us,
	}
}

func (s *Service) HandleGameAction(gameUUID uuid.UUID, userUUID uuid.UUID, guestToken string, action func(*domain.Game) error) (domain.Game, error) {
	game, err := s.gs.GetGame(gameUUID, userUUID, guestToken)
	if err != nil {
		return game, fmt.Errorf("game not found : %v", err)
	}

	err = action(&game)
	if err != nil {
		return domain.Game{}, err
	}
	err = s.gs.UpdateGame(game)
	return game, err
}

func (s *Service) GetGame(gameId uuid.UUID, userId uuid.UUID) (domain.Game, error) {
	game, err := s.gs.GetGame(gameId, userId, "")
	if err != nil {
		return domain.Game{}, err
	}
	if game.GameId == uuid.Nil() {
		return domain.Game{}, errors.New("no game found")
	}
	return game, nil
}

func (s *Service) GetUser(userUUID uuid.UUID) (domain.User, error) {
	user, err := s.us.GetUser(userUUID)
	if err != nil {
		return domain.User{}, err
	}

	if user == nil {
		user = &domain.User{
			UserId:  userUUID,
			Balance: defaultUserBalance,
		}

		if err := s.us.CreateUser(*user); err != nil {
			return domain.User{}, err
		}
	}

	return *user, nil
}

func (s *Service) GetUserAccount(userUUID uuid.UUID) (domain.UserAccountResponse, error) {
	user, err := s.GetUser(userUUID)
	if err != nil {
		return domain.UserAccountResponse{}, nil
	}
	games, err := s.GetGames(userUUID, 10, 0)
	if err != nil {
		return domain.UserAccountResponse{}, nil
	}
	return domain.UserAccountResponse{Balance: user.Balance, Games: games}, nil
}

func (s *Service) GetGames(userId uuid.UUID, limit, offset int) ([]domain.GameResponse, error) {
	games, err := s.gs.GetGames(userId, limit, offset)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.GameResponse, 0, len(games))
	for _, game := range games {
		summaries = append(summaries, game.ToResponse())
	}

	return summaries, nil
}

func (s *Service) GetNewGame(userID uuid.UUID) (domain.Game, error) {
	newGameUUID := uuid.New()
	newGame := domain.Game{GameId: newGameUUID, Deck: domain.GetNewDeck(domain.DeckConfig{Shuffle: true}), PlayerHand: domain.Hand{}, DealerHand: domain.Hand{}, Phase: "player"}

	if userID != uuid.Nil() {
		user, err := s.GetUser(userID)
		if err != nil {
			return domain.Game{}, err
		}
		newGame.UserId = user.UserId
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
	err := s.gs.UpdateGame(newGame)
	return newGame, err
}

func generateGuestToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
