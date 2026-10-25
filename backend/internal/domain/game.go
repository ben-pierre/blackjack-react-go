package domain

import (
	"errors"
	"math/rand/v2"
	"uuid"
)

type Card struct {
	Suit CardSuit `json:"suit"`
	Rank int      `json:"rank"`
}

type Hand struct {
	Cards []Card `json:"cards"`
}

type Game struct {
	UUID       uuid.UUID
	PlayerHand Hand
	DealerHand Hand
	Phase      string
	Deck       Deck
	Wager      int
}

type GameResponse struct {
	UUID       uuid.UUID `json:"id"`
	PlayerHand Hand      `json:"playerHand"`
	DealerHand Hand      `json:"dealerHand"`
	Phase      string    `json:"phase"`
}

type GameSummary struct {
	UUID  uuid.UUID `json:"id"`
	Phase string    `json:"phase"`
}

type Deck []Card

type CardSuit string

const (
	Diamonds CardSuit = "diamonds"
	Clubs    CardSuit = "clubs"
	Spades   CardSuit = "spades"
	Hearts   CardSuit = "hearts"
)

func (game *Game) ToResponse() GameResponse {
	dealerHand := game.DealerHand

	if game.Phase == "player" && len(dealerHand.Cards) > 1 {
		dealerHand.Cards = append([]Card(nil), dealerHand.Cards...)
		dealerHand.Cards[1] = Card{Suit: "X", Rank: 0}
	}

	return GameResponse{
		UUID:       game.UUID,
		PlayerHand: game.PlayerHand,
		DealerHand: dealerHand,
		Phase:      game.Phase,
	}
}

func (game *Game) Hit() error {
	if game.Phase != "player" {
		return errors.New("invalid game phase")
	}
	game.PlayerHand.Cards = append(game.PlayerHand.Cards, game.Deck.DrawCard())
	game.CheckWin()

	return nil
}

func (game *Game) Stay() error {
	if game.Phase != "player" {
		return errors.New("invalid game phase")
	}
	game.Phase = "dealer"
	return game.playDealerHand()
}

func (game *Game) CheckWin() {
	playerValue := game.PlayerHand.value()
	dealerValue := game.DealerHand.value()

	//Natural blackjacks
	if game.Phase == "player" && playerValue == 21 && len(game.PlayerHand.Cards) == 2 {
		if playerValue == dealerValue {
			game.Phase = "push"
			return
		}
		game.Phase = "win"
		return
	}

	// Player busts
	if playerValue > 21 {
		game.Phase = "lose"
		return
	}

	// Only resolve the dealer's hand after their turn begins
	if game.Phase != "dealer" {
		return
	}

	// Dealer busts
	if dealerValue > 21 {
		game.Phase = "win"
		return
	}

	// Compare final totals
	switch {
	case playerValue > dealerValue:
		game.Phase = "win"
	case playerValue < dealerValue:
		game.Phase = "lose"
	default:
		game.Phase = "push"
	}
}

func (game *Game) playDealerHand() error {
	if game.Phase != "dealer" {
		return errors.New("invalid game phase")
	}

	for game.DealerHand.value() < 17 {
		game.DealerHand.Cards = append(game.DealerHand.Cards, game.Deck.DrawCard())
	}
	game.CheckWin()
	return nil
}

func (hand Hand) value() int {
	numAces := 0
	total := 0
	for _, card := range hand.Cards {
		if card.Rank == 1 {
			total += 11
			numAces += 1
		} else {
			total += min(card.Rank, 10)
		}

	}
	for total > 21 && numAces > 0 {
		total -= 10
		numAces--
	}
	return total
}

func (d *Deck) DrawCard() Card {
	deck := *d
	card := deck[len(deck)-1]
	*d = deck[:len(deck)-1]
	return card
}

func GetNewDeck(cfg DeckConfig) []Card {
	suits := []CardSuit{Diamonds, Clubs, Spades, Hearts}
	deck := make([]Card, 0, 52)
	for _, suit := range suits {
		for rank := 1; rank <= 13; rank++ {
			deck = append(deck, Card{suit, rank})
		}
	}
	if cfg.Shuffle {
		rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	}
	return deck
}

type DeckConfig struct {
	Shuffle bool
}
