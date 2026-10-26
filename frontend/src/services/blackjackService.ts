import * as blackjack from "../types/blackjack";
import { fetchAuthSession } from "aws-amplify/auth";

class BlackJackService {
  private gameUUID: string = "";

  private readonly apiUrl = import.meta.env.VITE_API_URL;

  async newGame(): Promise<blackjack.Game> {

    const response = await apiFetchAndLog<blackjack.Game>(`${this.apiUrl}/blackjack/new`, {
      method: "POST",
      headers: await getAuthHeaders() 
    });
    if (!response.httpResponse.ok) {
      throw new Error("Failed to create game");
    }
    this.gameUUID = response.body.id;
    return response.body;
  }

  async hit(): Promise<blackjack.Game> {

    const response = await apiFetchAndLog<blackjack.Game>(
      `${this.apiUrl}/blackjack/${this.gameUUID}/hit`,
      { method: "POST",
      headers: await getAuthHeaders() 
    }
    );
    if (!response.httpResponse.ok) {
      throw new Error("Failed to create game");
    }
    return response.body;
  }

  async stay(): Promise<blackjack.Game> {
    const response = await apiFetchAndLog<blackjack.Game>(
      `${this.apiUrl}/blackjack/${this.gameUUID}/stay`,
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

  if (!token) {
    throw new Error("User is not authenticated");
  }

  return {
    Authorization: `Bearer ${token}`,
  };
}

export const blackjackService = new BlackJackService();
