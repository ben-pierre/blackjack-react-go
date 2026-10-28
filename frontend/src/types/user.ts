import type { Game } from "./blackjack"

export type User = {
    email : string,
    id : string,
    balance : number
}

export type UserAccount = {
    balance : number,
    games : Game[]
}