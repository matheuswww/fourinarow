import { handleWs } from "./handleWs";
import { basePath, wsbasePath } from "./path";

export async function createRoom(): Promise<WebSocket> {
  try {
    const response = await fetch(basePath+"/token", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Error fetching token: ${response.statusText}`);
    }

    const data: { token: string,  user_id: string } = await response.json();
    const token = data.token;
    window.localStorage.setItem("user_id", data.user_id)
    const ws = new WebSocket(wsbasePath+"/room");

    return handleWs(ws, token);
  } catch (error) {
    console.error("Error in request or WebSocket connection:", error);
    throw error;
  }
}
