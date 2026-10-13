
import "./Card.css";
import type { Card } from "../Types/PlayingCard";
import { CardSuit } from "../Types/PlayingCard";


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

  return (
    <div className={`card ${isRed ? "red" : "black"}`}>
      <span className="card-rank">{card.visible ? card.rank : "X"}</span>
      <span className="card-suit">{card.visible ? suitSymbols[card.suit] : "X"}</span>
    </div>
  );
}
