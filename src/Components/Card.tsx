
import "./Card.css";
import type { Card } from "../types/blackjack";
import { CardRank, CardSuit } from "../types/blackjack";


type CardProps = {
  card: Card;
};

export default function PlayingCard({ card }: CardProps) {
  const isRed = card.suit === CardSuit.Heart || card.suit === CardSuit.Diamond;

  const suitSymbols = {
    [CardSuit.Heart]: "♥",
    [CardSuit.Diamond]: "♦",
    [CardSuit.Club]: "♣",
    [CardSuit.Spade]: "♠"
  };

  function rankToChar(){
    switch(true){
      case(card.rank <=0): return "X"
      case(card.rank==CardRank.Ace):return "A"
      case(card.rank==CardRank.Jack): return "J"
      case(card.rank==CardRank.Queen): return "Q"
      case(card.rank==CardRank.King): return "K"
      default:return card.rank
    }
  }

  return (
    <div className={`card ${isRed ? "red" : "black"}`}>
      <span className="card-rank">{rankToChar()}</span>

      <span className="card-suit">{card.rank>0 ? suitSymbols[card.suit] : "X"}</span>
    </div>
  );
}
