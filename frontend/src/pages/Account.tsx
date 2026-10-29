import * as user from "../types/user";
import { useEffect, useState } from "react";
import { blackjackService } from "../services/blackjackService";
import GameHistoryCard from "../components/GameHistoryCard";

export default function Account() {
  const [userAccount, setUserAccount] = useState<
    user.UserAccount | undefined
  >();

  useEffect(() => {
    async function refreshUserAccount() {
      try {
        const user = await blackjackService.getAccount();
        setUserAccount(user);
      } catch {
        setUserAccount(undefined);
      }
    }
    refreshUserAccount();
  }, []);

  return (
    <>
      {userAccount ? (
        <div>
          <h1>Account</h1>

          <div className="accountBalance">
            <span className="accountBalanceLabel">Account Balance</span>
            <span className="accountBalanceAmount">
              ${userAccount.balance.toLocaleString()}
            </span>
          </div>

          <h2>Game History ({userAccount.games.length})</h2>

          <div className="gameHistory">
            {userAccount.games.map((game) => (
              <GameHistoryCard
                gameId={game.gameId}
                playerHand={game.playerHand}
                wager={game.wager}
                resultPhase={game.phase}
                dealerHand={game.dealerHand}
              />
            ))}
          </div>
        </div>
      ) : (
        <div>
          <span className="spinner" />
          <p>"Fetching account"</p>
        </div>
      )}
    </>
  );
}
