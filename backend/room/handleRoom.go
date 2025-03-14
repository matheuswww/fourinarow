package room

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func sendMessage(roomId string, msg1, msg2 any) {
	if msg1 != "" {
		rooms[roomId].user[0].conn.WriteJSON(msg1)
	}
	if msg2 != "" {
		rooms[roomId].user[1].conn.WriteJSON(msg2)
	}
}

func handleRoom(conn *websocket.Conn, roomdId, userId string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				if rooms[roomdId].user[0] != nil && rooms[roomdId].user[0].conn != nil && rooms[roomdId].user[1] != nil && rooms[roomdId].user[1].conn != nil {
					if rooms[roomdId].user[0].userId == userId {
						msg := "player disconnected"
						sendMessage(roomdId, "", msg)
						rooms[roomdId].user[0].conn = nil
					}
					if rooms[roomdId].user[1].userId == userId  {
						msg := "player disconnected"
						sendMessage(roomdId, msg, "")
						rooms[roomdId].user[1].conn = nil
					}
					break
				}
				delete(rooms, roomdId)
				log.Println("Room closed")
				break
			}
			log.Println("Error trying to read message: ", err)
			break
		}
		var play Request
		err = json.Unmarshal(message, &play)
		if err != nil {
			log.Println("Error trying Unmarshal")
			continue
		}
		if rooms[roomdId].user[0] != nil && rooms[roomdId].user[1] != nil {
			if rooms[roomdId].user[0].userId == userId {
				if rooms[roomdId].matrix[play.Play[0]][play.Play[1]] == 2 || !rooms[roomdId].user[0].play {
					msg := Message{ "invalid play" }
					sendMessage(roomdId, msg, msg)
					continue
				}
				rooms[roomdId].matrix[play.Play[0]][play.Play[1]] = 1
				rooms[roomdId].user[0].play = false
				rooms[roomdId].user[1].play = true
				msg := Response{ rooms[roomdId].matrix }
				sendMessage(roomdId, msg, msg)
			}
			if rooms[roomdId].user[1].userId == userId {
				if rooms[roomdId].matrix[play.Play[0]][play.Play[1]] == 1 || !rooms[roomdId].user[1].play {
					msg := Message{ "invalid play" }
					sendMessage(roomdId, msg, msg)
					continue
				}
				rooms[roomdId].matrix[play.Play[0]][play.Play[1]] = 2
				rooms[roomdId].user[1].play = false
				rooms[roomdId].user[0].play = true
				msg := Response{ rooms[roomdId].matrix }
				sendMessage(roomdId, msg, msg)
			}
		}
	}
}
