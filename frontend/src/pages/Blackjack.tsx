import { useState } from "react";
import * as blackjack from "../types/blackjack";
import { blackjackService } from "../services/blackjackService";
import Table from "../components/Table";

function getEmptyGame(): blackjack.Game {
  const game: blackjack.Game = {
    gameId: "",
    playerHand: { cards: [] },
    dealerHand: { cards: [] },
    phase: "preGame",
    guestToken: "",
    wager: 0,
  };
  return game;
}

export default function Blackjack() {
  const [game, setGame] = useState<blackjack.Game>(getEmptyGame());

  const [visibleDealerCards, setVisibleDealerCards] =
    useState<blackjack.playingHand>({ cards: [] });

  const [dealerRevealed, setDealerRevealed] = useState(false);

  const gameOver =
    game.phase === "win" || game.phase === "lose" || game.phase === "push";

  async function newGame() {
    setGameUi(await blackjackService.newGame());
  }

  function delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  async function hit() {
    if (game.phase !== "player") return;
    setGameUi(await blackjackService.hit());
  }

  function setGameUi(game: blackjack.Game) {
    setGame(game);
    setVisibleDealerCards({ cards: game.dealerHand.cards });
  }

  async function stay() {
    if (game.phase !== "player") return;
    const updatedGame = await blackjackService.stay();
    setDealerRevealed(true);
    await delay(500);
    const cards = updatedGame.dealerHand.cards;

    setVisibleDealerCards({ cards: cards.slice(0, 2) });

    for (let index = 2; index < cards.length; index++) {
      await delay(500);
      setVisibleDealerCards((prev) => ({
        cards: [...prev.cards, cards[index]],
      }));
    }
    setGameUi(updatedGame);
  }

  function getPhaseTitle(): string {
    switch (game.phase) {
      case "player":
        return "Your Turn";
      case "player":
        return "Dealer's Turn";
      case "win":
        return "You Win!";
      case "lose":
        return "You Lose";
      case "push":
        return "Push!";
      default:
        return "Waiting...";
    }
  }

  function getResultMessage(): string {
    switch (game.phase) {
      case "win":
        return "Nice hand! You beat the dealer.";
      case "lose":
        return blackjack.sumHand(game.playerHand) > 21
          ? "Bust! You went over 21."
          : "Better luck next hand.";
      case "push":
        return "A push! Nobody wins this hand.";
      default:
        return "Waiting...";
    }
  }

  return (
    <main className="blackjack">
      <header className="blackjackHeader">
        <h1>Blackjack</h1>
      </header>

      {
        <Table
          visibleDealerCards={visibleDealerCards}
          dealerRevealed={dealerRevealed}
          playerHand={game?.playerHand}
        />
      }

      {gameOver ? (
        <section
          className={`gameResult result-${game.phase}`}
          aria-live="polite"
        >
          <div className="resultIcon">
            {game.phase === "win" && "♠"}
            {game.phase === "lose" && "♣"}
            {game.phase === "push" && "♦"}
          </div>

          <h2>{getPhaseTitle()}</h2>
          <p>{getResultMessage()}</p>

          <div className="gameActions">
            <button type="button" onClick={newGame}>
              Play Again
            </button>
          </div>
        </section>
      ) : game.phase === "preGame" ? (
        <div className="gameActions">
          <button type="button" onClick={newGame}>
            Deal
          </button>
        </div>
      ) : (
        <div className="gameActions">
          <button
            type="button"
            onClick={hit}
            disabled={game.phase !== "player"}
          >
            Hit
          </button>
          <button
            type="button"
            onClick={stay}
            disabled={game.phase !== "player"}
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
