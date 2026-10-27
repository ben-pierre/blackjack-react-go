import * as blackjack from "../types/blackjack";
import { fetchAuthSession } from "aws-amplify/auth";

class BlackJackService {
  private game: blackjack.Game

  private readonly apiUrl = import.meta.env.VITE_API_URL;

  async newGame(): Promise<blackjack.Game> {

    const response = await apiFetchAndLog<blackjack.Game>(`${this.apiUrl}/blackjack/new`, {
      method: "POST",
      headers: await getAuthHeaders() 
    });
    if (!response.httpResponse.ok) {
      throw new Error("Failed to create game");
    }
    this.game = response.body
    return response.body;
  }

  async hit(): Promise<blackjack.Game> {

    const response = await apiFetchAndLog<blackjack.Game>(
      `${this.apiUrl}/blackjack/${this.game.gameId}/hit`,
      { method: "POST",
      headers: {...await getAuthHeaders(), ...getGuestTokenHeader(this.game.guestToken)}
    }
    );
    if (!response.httpResponse.ok) {
      throw new Error("Failed to create game");
    }
    return response.body;
  }

  async stay(): Promise<blackjack.Game> {
    const response = await apiFetchAndLog<blackjack.Game>(
      `${this.apiUrl}/blackjack/${this.game.gameId}/stay`,
      { method: "POST",
      headers: await getAuthHeaders()  }
    );
    if (!response.httpResponse.ok) {
      throw new Error("Failed to create game");
    }
    return response.body;
  }
}

async function apiFetchAndLog<T>(url: string, options?: RequestInit) : Promise<{httpResponse : Response, body:T}> {
  console.log("REQUEST:", options?.method ?? "GET", url);
  const response = await fetch(url, options);
  const body = await response.json();
  console.log("RESPONSE:", response.status, JSON.stringify(body, null, 2));
  return { httpResponse: response, body };
}

async function getAuthHeaders(): Promise<HeadersInit> {
  const session = await fetchAuthSession();
  const token = session.tokens?.accessToken.toString();

  if(!token){
    return
  }

  return {
    Authorization: `Bearer ${token}`,
  };
}

function getGuestTokenHeader(guestToken: string): HeadersInit{

  if(guestToken!==""){
    return {"X-Guest-Token": guestToken}
  }else{
    return {}
  }
}

export const blackjackService = new BlackJackService();
