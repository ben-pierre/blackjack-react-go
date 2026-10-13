export type Card = {
    rank: CardRank;
    suit: CardSuit;
    visible: boolean;
  };


export enum CardSuit{
    Spade,
    Diamond,
    Club,
    Heart
}

export enum CardRank {
    Ace = 1,
    Two,
    Three,
    Four,
    Five,
    Six,
    Seven,
    Eight,
    Nine,
    Ten,
    Jack =10,
    Queen = 10,
    King = 10
  }

