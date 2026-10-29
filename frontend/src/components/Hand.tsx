import PlayingCard from "./Card";
import { sumHand } from "../types/blackjack";
import type { playingHand } from "../types/blackjack";

type HandProps = {
  title: string;
  hand: playingHand;
  hidden?: boolean;
  wager? : number
};

export default function Hand({
  title,
  hand,
  hidden = false,
  wager = 0
}: HandProps) {
  return (
    <div className="handSection">
      <div className="handHeader">
        <h2>{title}</h2>

        <span className="handValue">
          {hidden ? "?" : sumHand(hand)}
        </span>

        {wager > 0 && <span className="handWager">Wager: ${wager}</span>}

      </div>

      <div className="handContainer">
        {hand.cards.map((card, index) => (
          <PlayingCard key={index} card={card} />
        ))}
      </div>
    </div>
  );
}