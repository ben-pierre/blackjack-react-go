import Hand from "./Hand"
import type { playingHand } from "../types/blackjack"

export type TableProperties = {
    visibleDealerCards : playingHand
    playerHand : playingHand
    dealerRevealed : boolean
}

export default function Table({visibleDealerCards, playerHand, dealerRevealed}:TableProperties){

    return(
        <section className="table">
        <Hand
          title="Dealer"
          hand={visibleDealerCards}
          hidden={!dealerRevealed}
        />

        <div className="tableDivider" />
        <Hand title="Your hand" hand={playerHand}/>
      </section>


    )
}