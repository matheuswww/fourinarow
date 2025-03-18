interface Play {
  play: [number, number];
}

function sendPlay(ws: WebSocket, move: [number, number]): void {
  if (ws.readyState === WebSocket.OPEN) {
    const payload: Play = { play: move };
    ws.send(JSON.stringify(payload));
  } else {
    console.error('WebSocket is not open');
  }
}
