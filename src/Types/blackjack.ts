import { CardSuit, CardRank } from "../Types/PlayingCard";
import type { Card } from "../Types/PlayingCard";

export type playingHand = {
  cards: Card[];
};

export type GamePhase =
  | "playerPlaying"
  | "dealerPlaying"
  | "win"
  | "lose"
  | "tie";

export function sumHand(hand: playingHand): number {
  let total = 0;
  let aces = 0;
  for (const card of hand.cards) {
    if (card.rank === CardRank.Ace) {
      aces++;
      total += 11;
    } else {
      total += card.rank;
    }
  }
  while (total > 21 && aces > 0) {
    total -= 10;
    aces--;
  }
  return total;
}

export function checkDealerWin(
  dealerHandValue: number,
  playerHandValue: number
): GamePhase {
    console.log(`Checking dealer win for ${dealerHandValue} d vs ${playerHandValue} p`)
  if (dealerHandValue > 21) {
    return "win";
  }
  if (dealerHandValue > playerHandValue && dealerHandValue >= 17) {
    return "lose";
  }
  if (dealerHandValue === playerHandValue && dealerHandValue >= 17) {
    return "tie";
  }
  if (dealerHandValue >= 17) {
    return "win";
  }
  return "dealerPlaying"; //default
}

export function createDeck(): Card[] {
  const deck: Card[] = [];

  for (const suit of Object.values(CardSuit).filter(
    (v): v is CardSuit => typeof v === "number"
  )) {
    for (const rank of Object.values(CardRank).filter(
      (v): v is CardRank => typeof v === "number"
    )) {
      deck.push({ suit, rank, visible: true });
    }
  }

  for (let i = deck.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [deck[i], deck[j]] = [deck[j], deck[i]];
  }

  return deck;
}