package room

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func sendMessage(roomId string, msg1, msg2 string) {
	if msg1 != "" {
		rooms[roomId][0].conn.WriteJSON(Message{ msg1 })
	}
	if msg2 != "" {
		rooms[roomId][1].conn.WriteJSON(Message{ msg1 })
	}
}

func handleRoom(conn *websocket.Conn, id string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				if 	rooms[id][0] != nil || rooms[id][1] != nil {
					if rooms[id][0].conn != nil {
						msg := "player disconnected"
						sendMessage(id, msg, "")
						rooms[id][1].conn = nil
					}
					if rooms[id][1].conn != nil {
						msg := "player disconnected"
						sendMessage(id, "", msg)
						rooms[id][0].conn = nil
					}
					break
				}
				delete(rooms, id)
				log.Println("Room closed")
				break
			}
			log.Println("Error trying to read message: ", err)
			break
		}
		var msg Request
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Error trying Unmarshal")
			continue
		}
	}
}
