export type playingHand = {
  cards: Card[];
};


export type Game = {
  gameId : string
  guestToken : string
  playerHand : playingHand
  dealerHand : playingHand
  phase : GamePhase
}

export type GamePhase =
  | "preGame"
  | "player"
  | "dealer"
  | "win"
  | "lose"
  | "push";

export function sumHand(hand: playingHand): number {

  let total = 0;
  let aces = 0;
  for (const card of hand.cards) {
    if (card.rank === CardRank.Ace) {
      aces++;
      total += 11;
    } else {
      total += Math.min(card.rank,10);
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
    return "push";
  }
  if (dealerHandValue >= 17) {
    return "win";
  }
  return "dealer"; //default
}

export function createDeck(): Card[] {
  const deck: Card[] = [];

  for (const suit of Object.values(CardSuit).filter(
    (v): v is CardSuit => typeof v === "number"
  )) {
    for (const rank of Object.values(CardRank).filter(
      (v): v is CardRank => typeof v === "number"
    )) {
      deck.push({ suit, rank});
    }
  }

  for (let i = deck.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [deck[i], deck[j]] = [deck[j], deck[i]];
  }

  return deck;
}
export type Card = {
  rank: CardRank;
  suit: CardSuit;
};

export const CardSuit = {
  Spade: "spades",
  Diamond: "diamonds",
  Club: "clubs",
  Heart:"hearts",
} as const;

export type CardSuit = typeof CardSuit[keyof typeof CardSuit];

export const CardRank = {
  Ace: 1,
  Two: 2,
  Three: 3,
  Four: 4,
  Five: 5,
  Six: 6,
  Seven: 7,
  Eight: 8,
  Nine: 9,
  Ten: 10,
  Jack: 11,
  Queen: 12,
  King: 13,
} as const;

export type CardRank = typeof CardRank[keyof typeof CardRank];