import "./Blackjack.css";
import { useState } from "react";
import { useNavigate, Link} from "react-router";
import { checkDealerWin, sumHand, createDeck } from "../Types/blackjack";
import HandSection from "../Components/Hand";
import type { playingHand, GamePhase } from "../Types/blackjack";
import type { Card } from "../Types/PlayingCard";

function dealInitialHands(deck: Card[]) {
  const newDeck = [...deck];

  const playerHand: playingHand = {
    cards: [newDeck.pop()!, newDeck.pop()!],
  };

  const dealerHand: playingHand = {
    cards: [newDeck.pop()!, newDeck.pop()!],
  };

  dealerHand.cards[0].visible = false;

  return { newDeck, playerHand, dealerHand };
}

export default function Blackjack() {
  const [initialGame] = useState(() => dealInitialHands(createDeck()));

  const [phase, setPhase] = useState<GamePhase>("playerPlaying");
  const [deck, setDeck] = useState<Card[]>(initialGame.newDeck);
  const [playerHand, setPlayerHand] = useState<playingHand>(
    initialGame.playerHand
  );
  const [dealerHand, setDealerHand] = useState<playingHand>(
    initialGame.dealerHand
  );
  const [dealerRevealed, setDealerRevealed] = useState(false);

  const navigate = useNavigate();

  const gameOver = phase === "win" || phase === "lose" || phase === "tie";

  function resetGame() {
    const game = dealInitialHands(createDeck());

    setDealerRevealed(false);
    setDeck(game.newDeck);
    setPlayerHand(game.playerHand);
    setDealerHand(game.dealerHand);
    setPhase("playerPlaying");
  }

  function delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  function showHand(hand: playingHand): playingHand {
    return {
      cards: hand.cards.map((card) => ({
        ...card,
        visible: true,
      })),
    };
  }

  function hit() {
    if (phase !== "playerPlaying" || deck.length === 0) return;

    const nextDeck = [...deck];
    const card = nextDeck.pop()!;

    const nextHand: playingHand = {
      cards: [...playerHand.cards, card],
    };

    setDeck(nextDeck);
    setPlayerHand(nextHand);

    if (sumHand(nextHand) > 21) {
      setDealerHand(showHand(dealerHand));
      setPhase("lose");
    }
  }

  async function stay() {
    if (phase !== "playerPlaying") return;
    
    setPhase("dealerPlaying");

    // Pause before revealing the dealer's hole card.
    await delay(1000);

    let nextDealerHand = showHand(dealerHand);
    let nextDeck = [...deck];

    setDealerHand(nextDealerHand);
    setDealerRevealed(true);

    // Draw one card at a time, with a pause between each.
    while (sumHand(nextDealerHand) < 17 && nextDeck.length > 0) {
      await delay(1000);

      const card = nextDeck.pop()!;

      nextDealerHand = {
        cards: [...nextDealerHand.cards, card],
      };

      setDealerHand(nextDealerHand);
    }

    setDeck(nextDeck);

    // Determine the result after the dealer finishes.
    setPhase(checkDealerWin(sumHand(nextDealerHand), sumHand(playerHand)));
  }

  function getPhaseTitle(): string {
    switch (phase) {
      case "playerPlaying":
        return "Your Turn";
      case "dealerPlaying":
        return "Dealer's Turn";
      case "win":
        return "You Win!";
      case "lose":
        return "You Lose";
      case "tie":
        return "Push!";
    }
  }

  function getResultMessage(): string {
    switch (phase) {
      case "win":
        return "Nice hand! You beat the dealer.";
      case "lose":
        return sumHand(playerHand) > 21
          ? "Bust! You went over 21."
          : "Better luck next hand.";
      case "tie":
        return "A push! Nobody wins this hand.";
      default:
        return "";
    }
  }

  return (
    <main className="blackjack">
      <header className="blackjackHeader">
      <Link to="/" className="homeButton">
      ← Home
</Link>


        <h1>Blackjack</h1>
      </header>

      <section className="table">
      <HandSection
          title="Dealer" hand={dealerHand} hidden={!dealerRevealed}
        />

        <div className="tableDivider" />
        <HandSection
          title="Your hand" hand={playerHand}
        />
      </section>

      {gameOver ? (
        <section className={`gameResult result-${phase}`} aria-live="polite">
          <div className="resultIcon">
            {phase === "win" && "♠"}
            {phase === "lose" && "♣"}
            {phase === "tie" && "♦"}
          </div>

          <h2>{getPhaseTitle()}</h2>
          <p>{getResultMessage()}</p>

          <div className="gameActions">
            <button type="button" onClick={resetGame}>
              Play Again
            </button>
            <Link to="/" className="secondaryButton">
  Home
</Link>
          </div>
        </section>
      ) : (
        <div className="gameActions">
          <button
            type="button"
            onClick={hit}
            disabled={phase !== "playerPlaying"}
          >
            Hit
          </button>

          <button
            type="button"
            onClick={stay}
            disabled={phase !== "playerPlaying"}
          >
            Stay
          </button>
        </div>
      )}

      <footer className="blackjackFooter">
        Dealer stands on 17 · Good luck!
      </footer>
    </main>
  );
}
