import "./Blackjack.css";
import { useState, useRef, useEffect } from "react";
import PlayingCard from "../Components/Card";
import EndGameModal from "../Components/blackjack/EndGameModal";
import { checkDealerWin, sumHand, createDeck } from "../Types/blackjack";
import type { playingHand, GamePhase } from "../Types/blackjack";
import type { Card } from "../Types/PlayingCard";
import { useNavigate } from "react-router";


export default function Blackjack() {

  const [phase, setPhase] = useState<GamePhase>("playerPlaying");

  const [deck, setDeck] = useState<Card[]>(createDeck());

  const [playerHand, setPlayerHand] = useState<playingHand>({
    cards: [],
  });

  const [dealerHand, setDealerHand] = useState<playingHand>({
    cards: [],
  });

  const [,] = useState(() => resetGame());

  const navigate = useNavigate();

  function goHome() {
    navigate("/");
  }


  function resetGame() {
    changePhase("playerPlaying");

    const newDeck = createDeck();
    setPlayerHand({ cards: [newDeck.pop()!, newDeck.pop()!] });
    const dealerCards = [newDeck.pop()!, newDeck.pop()!];
    dealerCards[0].visible = false;
    setDealerHand({ cards: dealerCards });

    setDeck(newDeck);
  }

  function hit() {
    if (phase === "playerPlaying") drawPlayerCard();
    if (phase === "dealerPlaying") drawDealerCard();
  }

  function stay() {
    if (phase !== "playerPlaying") return;
    startDealerPhase();
  }

  function startDealerPhase() {
    changePhase("dealerPlaying");
    setDealerHand(showHand(dealerHand));
    changePhase(checkDealerWin(sumHand(dealerHand), sumHand(playerHand)));
  }

  function drawPlayerCard() {
    const nextHand = drawCard(playerHand);
    setPlayerHand(nextHand);

    if (sumHand(nextHand) > 21) {
      changePhase("lose");
    }
  }

  function showHand(hand: playingHand): playingHand {
    return {
      cards: hand.cards.map((card) => ({
        ...card,
        visible: true,
      })),
    };
  }

  function drawDealerCard() {
    const nextHand = drawCard(dealerHand);
    setDealerHand(nextHand);
    changePhase(checkDealerWin(sumHand(nextHand), sumHand(playerHand)));
  }

  function changePhase(nextPhase: GamePhase) {
    setPhase(nextPhase);
  }

  function getPhaseTitle(): string {
    switch (phase) {
      case "playerPlaying":
        return "Your Turn";
      case "dealerPlaying":
        return "Dealer's Turn";
      case "win":
        return "You won!";
      case "lose":
        return "You lost!";
      case "tie":
        return "You tied!";
      default:
        return "";
    }
  }

  function drawCard(hand: playingHand): playingHand {
    const nextDeck = [...deck];
    const card = nextDeck.pop();
    setDeck(nextDeck);
    return {
      cards: [...hand.cards, card],
    };
  }

  return (
    <>

      <EndGameModal message={getPhaseTitle()} open = {(phase === "win" || phase === "lose" || phase === "tie")} goHome={goHome} resetGame={resetGame}/>

      <div>
        <h1>Blackjack</h1>
        <h2>{getPhaseTitle()}</h2>
        <button type="button" onClick={hit}>
          Hit
        </button>
        <button type="button" onClick={stay}>
          Stay
        </button>
        <h2>Your Hand</h2>
        <div className="handContainer">
          {playerHand.cards.map((handCard, index) => (
            <PlayingCard key={index} card={handCard}/>
          ))}
        </div>
        <h3>Value</h3>
        <p>{sumHand(playerHand).toString()}</p>
        <h2>Dealer Hand</h2>
        <div className="handContainer">
          {dealerHand.cards.map((handCard, index) => (
            <PlayingCard key={index} card={handCard} />
          ))}
        </div>
        <h3>Value</h3>
        <p>
          {phase === "playerPlaying" ? "X" : sumHand(dealerHand).toString()}
        </p>
      </div>
    </>
  );
}
