import { handleWs } from "./handleWs";
import { wsbasePath } from "./path";

export async function connectRoom(roomId: string): Promise<WebSocket> {
  try {
    let token = window.localStorage.getItem("token")
    const ws = new WebSocket(wsbasePath+`/room?id=${roomId}`);
    if (!token) {
      token = ""
    }
    return handleWs(ws, token);
  } catch (error) {
    console.error("Error in request or WebSocket connection:", error);
    throw error;
  }
}