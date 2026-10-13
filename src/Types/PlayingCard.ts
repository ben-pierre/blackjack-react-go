export type Card = {
  rank: CardRank;
  suit: CardSuit;
  visible: boolean;
};

export const CardSuit = {
  Spade: 0,
  Diamond: 1,
  Club: 2,
  Heart: 3,
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
  Jack: 10,
  Queen: 10,
  King: 10,
} as const;

export type CardRank = typeof CardRank[keyof typeof CardRank];