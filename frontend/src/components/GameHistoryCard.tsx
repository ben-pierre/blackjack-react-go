import type { playingHand, GamePhase} from "../types/blackjack";
import { Link } from "react-router";
import Hand from "./Hand"

export type GameHistoryCardProperties = {
    gameId:string
    dealerHand: playingHand;
    playerHand: playingHand;
    wager: number;
    resultPhase: GamePhase;
  };
  
  export default function GameHistoryCard({
    gameId,
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
            return "You left the table! Resume?";
        }
      };


    return (
      <section className={"gameHistoryCard"}>
        <div className={`historyHeader result-${resultPhase}`}>
          {resultPhase === "player" ? (

<h3><Link to={`/blackjack/${gameId}`}>{getTitle()} </Link></h3>
          ):(
            <h3>{getTitle()}</h3>

          )}
          {wager > 0 && <span>Wager: ${wager}</span>}

        </div>
  
        <Hand
          title="Dealer"
          hand={dealerHand}
          hidden={resultPhase==="player"}
        />
  
        <div className="historyDivider" />
  
        <Hand
          title="Your hand"
          hand={playerHand}
        />
      </section>
    );
  }