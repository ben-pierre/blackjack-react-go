import * as blackjack from "../types/blackjack";

class BlackJackService {
  private gameUUID: string = "";

  private readonly apiUrl = import.meta.env.VITE_API_URL;

  async newGame(): Promise<blackjack.Game> {
    const response = await apiFetch<blackjack.Game>(`${this.apiUrl}/blackjack/new`, {
      method: "POST",
    });
    if (!response.httpResponse.ok) {
      throw new Error("Failed to create game");
    }
    this.gameUUID = response.body.id;
    return response.body;
  }

  async hit(): Promise<blackjack.Game> {
    const response = await apiFetch<blackjack.Game>(
      `${this.apiUrl}/blackjack/${this.gameUUID}/hit`,
      { method: "POST" }
    );
    if (!response.httpResponse.ok) {
      throw new Error("Failed to create game");
    }
    return response.body;
  }
  async stay(): Promise<blackjack.Game> {
    const response = await apiFetch<blackjack.Game>(
      `${this.apiUrl}/blackjack/${this.gameUUID}/stay`,
      { method: "POST" }
    );
    if (!response.httpResponse.ok) {
      throw new Error("Failed to create game");
    }
    return response.body;
  }
}

async function apiFetch<T>(url: string, options?: RequestInit) : Promise<{httpResponse : Response, body:T}> {
  console.log("REQUEST:", options?.method ?? "GET", url);
  const response = await fetch(url, options);
  const body = await response.json();
  console.log("RESPONSE:", response.status, JSON.stringify(body, null, 2));
  return { httpResponse: response, body };
}

export const blackjackService = new BlackJackService();
