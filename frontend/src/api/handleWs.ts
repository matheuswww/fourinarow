export function handleWs(ws: WebSocket, token: string): Promise<WebSocket> {
  return new Promise((resolve, reject) => {
    ws.onopen = () => {
      console.log("Connected to WebSocket successfully!");

      ws.send(JSON.stringify({ token: token }));
      console.log("Token sent to WebSocket:", token);

      resolve(ws);
    };

    ws.onerror = (error: Event) => {
      console.error("WebSocket error:", error);
      reject(error);
    };

    ws.onclose = () => {
      console.log("WebSocket connection closed.");
    };
  });
}