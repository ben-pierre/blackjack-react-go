import PlayingCard from "./Card";
import { sumHand } from "../types/blackjack";
import type { playingHand } from "../types/blackjack";

type HandProps = {
  title: string;
  hand: playingHand;
  hidden?: boolean;
};

export default function Hand({
  title,
  hand,
  hidden = false,
}: HandProps) {
  return (
    <div className="handSection">
      <div className="handHeader">
        <h2>{title}</h2>

        <span className="handValue">
          {hidden ? "?" : sumHand(hand)}
        </span>
      </div>

      <div className="handContainer">
        {hand.cards.map((card, index) => (
          <PlayingCard key={index} card={card} />
        ))}
      </div>
    </div>
  );
}