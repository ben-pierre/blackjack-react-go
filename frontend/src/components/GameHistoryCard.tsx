import type { playingHand, GamePhase} from "../types/blackjack";
import Hand from "./Hand"

export type GameHistoryCardProperties = {
    dealerHand: playingHand;
    playerHand: playingHand;
    wager: number;
    resultPhase: GamePhase;
  };
  
  export default function GameHistoryCard({
    dealerHand,
    playerHand,
    wager,
    resultPhase
  }: GameHistoryCardProperties) {

    const getTitle = () => {
        switch (resultPhase) {
          case "win":
            return "You won!";
          case "push":
            return "Push!";
          case "lose":
            return "You lost!";
          default:
            return "You left the table!";
        }
      };


    return (
      <section className={"gameHistoryCard"}>
        <div className={`historyHeader result-${resultPhase}`}>
          <h3>{getTitle()}</h3>
          <span>Wager: ${wager}</span>
        </div>
  
        <Hand
          title="Dealer"
          hand={dealerHand}
        />
  
        <div className="historyDivider" />
  
        <Hand
          title="Your hand"
          hand={playerHand}
        />
      </section>
    );
  }