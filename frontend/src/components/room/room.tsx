"use client";
import { useEffect, useState } from "react";
import styles from "./room.module.css";
import Link from "next/link";

interface RoomProps {
  ws: WebSocket | null;
  owner: boolean
}

export default function Room({ ws,owner }: RoomProps) {
  const [grid, setGrid] = useState<number[]>(Array(6 * 7).fill(0));
  const [start, setStart] = useState<boolean>(false)
  const [count, setCount] = useState<number>(-1)
  const [timeToPlay, setTimeToPlay] = useState<number>(0)
  const [winner, setWinner] = useState<"vermelha" | "amarela" | null>(null)

  function getPlayableRow(column: number) {
    for (let row = 5; row >= 0; row--) {
      const index = row * 7 + column;
      if (grid[index] === 0) return row;
    }
    return -1;
  }

  function clearColumnStyles(column: number) {
    const container = document.querySelector(`.${styles.container}`);
    if (column >= 0 && column < 7) {
      const cells = container?.querySelectorAll(`[data-column="${column}"]`);
      cells?.forEach((cell) => {
        cell.classList.remove(styles.playable);
      });
    }
  }

  function applyColumnStyles(column: number) {
    const container = document.querySelector(`.${styles.container}`);
    const cells = container?.querySelectorAll(`[data-column="${column}"]`);
    const playableRow = getPlayableRow(column);
    if (playableRow !== -1) {
      const playableCell = cells?.[playableRow];
      if (playableCell) {
        playableCell.classList.add(styles.playable);
      }
    }
  }

  function handleMouseMove(event: React.MouseEvent<HTMLDivElement>) {
    const target = event.target as HTMLElement;
    if (target.classList.contains(styles.cell)) {
      const column = Number(target.getAttribute("data-column"));
      const prevColumn = Number(
        event.currentTarget.style.getPropertyValue("--hovered-column") || -1
      );

      if (prevColumn !== column && prevColumn >= 0 && prevColumn < 7) {
        clearColumnStyles(prevColumn);
      }
      if (column >= 0 && column < 7) {
        event.currentTarget.style.setProperty("--hovered-column", String(column));
        applyColumnStyles(column);
      }
    }
  }

  function handleMouseLeave() {
    const container = document.querySelector(`.${styles.container}`);
    if (container instanceof HTMLElement) {
      const prevColumn = Number(container.style.getPropertyValue("--hovered-column") || -1);
      container.style.setProperty("--hovered-column", "-1");
      if (prevColumn >= 0 && prevColumn < 7) {
        clearColumnStyles(prevColumn);
      }
    }
  }

  function sendPlay(ws: WebSocket, move: [number, number]): void {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ play: move }));
    } else {
      console.error("WebSocket is not open");
    }
  }

  function handleCellClick(column: number) {
    const playableRow = getPlayableRow(column);
    if (playableRow !== -1 && ws) {
      sendPlay(ws, [playableRow, column]);
      clearColumnStyles(column);
    }
  }

  useEffect(() => {
    if (ws) {
      const userId = window.localStorage.getItem("user_id")
      let hasWinner = false
      ws.onclose = () => {
        window.location.href = "/"
      };
      ws.onmessage = (event: MessageEvent) => {
        try {
        const data = JSON.parse(event.data);
        console.log(data);
        if (Array.isArray(data.matrix)) {
          const newMatrix = data.matrix.flat();
          setGrid(newMatrix)
        }
        
        if (data.message && typeof data.message == "string" ) {
          if(data.message.startsWith("Time to play: ")) {
            const init = data.message.indexOf("Time to play: ") + "Time to play: ".length;
            const id = data.message.substring(init)
            if ((userId == id && owner) || (userId != id && !owner)) {
              setTimeToPlay(2)
            } else {
              setTimeToPlay(1)
            }
            return
          }
          if(data.message.startsWith("Timer: ")) {
            const init = data.message.indexOf("Timer: ") + "Timer: ".length;
            const time = Number(data.message.substring(init))
            if (!isNaN(time)){
              if (time == 0) {
                setTimeToPlay((t) => t - 1 == 0 ? 2 : 1)
              }
              setCount(time)
            }
            return
          }
          if(data.message.startsWith("Time to start: ")) {
            const init = data.message.indexOf("Time to start: ") + "Time to start: ".length;
            const time = Number(data.message.substring(init))
            if (!isNaN(time)){
              if (time == 0) {
                setStart(true)
                return
              }
              setCount(time)
            }
            return
          }
          
          if(data.message.startsWith("Winner: ")) {
            hasWinner = true
            const init = data.message.indexOf("Winner: ") + "Winner: ".length;
            const winnerId = data.message.substring(init)
            if ((userId == winnerId && owner) || (userId != winnerId && !owner)) {
              setWinner("vermelha")
            } else {
              setWinner("amarela")
            }
            setCount(-1)
            return
          }
          if(data.message.startsWith("Player disconnected: ") && !hasWinner) {
            window.location.href = "/"
            return
          }
        }
        } catch (error) {
          console.error("Error parsing message:", error);
        }
      }
    }
  }, [ws]);

  return (
    <div className={styles.page}>
      <div className={`${styles.header} ${(winner || !start) && styles.lowOpacity}`}>
        <h1>Timer: {count >= 0 ? count : 0}/5</h1>
        <span className={`${styles.play} ${ timeToPlay == 1 && styles.filledPlayer1} ${timeToPlay == 2 && styles.filledPlayer2}`}></span>
      </div>
      <div
          className={`${styles.container} ${(winner || !start) && styles.lowOpacity}`}
          onMouseMove={handleMouseMove}
          onMouseLeave={handleMouseLeave}
        >
          {Array.from({ length: 6 * 7 }).map((_, index) => {
            const column = index % 7;
            const row = Math.floor(index / 7);
            return (
              <div
                key={index}
                className={`${styles.cell} ${grid[index] === 1
                    ? styles.filledPlayer1
                    : grid[index] === 2
                      ? styles.filledPlayer2
                      : ""}`}
                data-column={column}
                data-row={row}
                onClick={() => handleCellClick(column)} />
            );
          })}
        </div>
        {(winner != null || !start) && <div className={styles.box}>
            { winner != null && 
              <>
                <p>{winner == "vermelha" ? "Vermelha" : winner == "amarela" && "Amarela"} Venceu!</p>
                <Link href={"/"} onClick={(e) => {e.preventDefault();window.location.href = "/"}}>Voltar para home</Link> 
              </>
            }
            { !winner && ( count == -1 ? <p>Aguardando um jogador...</p> : <p>Começa em: {count} <br /> Sua peça é: {owner ? "vermelha" : "amarela"}</p> ) }
          </div>
      }
    </div>
  );
}