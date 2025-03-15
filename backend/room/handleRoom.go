package room

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var timeToPlay = 5
var timStart = 3
var timeExit = 10

func sendMessage(roomId string, msg1, msg2 any) {
	room, ok := rooms[roomId]
	if !ok {
		return
	}
	if msg1 != "" && room.user[0] != nil && room.user[0].conn != nil {
		room.user[0].conn.WriteJSON(msg1)
	}
	if msg2 != "" && room.user[1] != nil && room.user[1].conn != nil {
		room.user[1].conn.WriteJSON(msg2)
	}
}

func handleRoom(conn *websocket.Conn, roomdId, userId string) {
	room := rooms[roomdId]
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				if room.user[0] != nil && room.user[0].conn != nil && room.user[1] != nil && room.user[1].conn != nil {
					if room.user[0].userId == userId {
						msg := "player disconnected"
						sendMessage(roomdId, "", msg)
						room.user[0].conn = nil
						room.user[0].exit = make(chan bool)
						handleTimeExit(room.user[0])
					}
					if room.user[1].userId == userId {
						msg := "player disconnected"
						sendMessage(roomdId, msg, "")
						room.user[1].conn = nil
						room.user[1].exit = make(chan bool)
						handleTimeExit(room.user[1])
					}
					break
				}
				rooms[roomdId].timer <- false
				close(rooms[roomdId].timer)
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
		if room.user[0] != nil && room.user[1] != nil {
			if room.user[0].userId == userId {
				if rooms[roomdId].matrix[play.Play[0]][play.Play[1]] == 2 || rooms[roomdId].play != 1 {
					msg := Message{"invalid play"}
					sendMessage(roomdId, msg, msg)
					continue
				}
				rooms[roomdId].matrix[play.Play[0]][play.Play[1]] = 1
				rooms[roomdId].play = 2
				msg := Response{rooms[roomdId].matrix}
				sendMessage(roomdId, msg, msg)
			}
			if room.user[1].userId == userId {
				if rooms[roomdId].matrix[play.Play[0]][play.Play[1]] == 1 || rooms[roomdId].play != 2 {
					msg := Message{"invalid play"}
					sendMessage(roomdId, msg, msg)
					continue
				}
				rooms[roomdId].matrix[play.Play[0]][play.Play[1]] = 2
				rooms[roomdId].play = 1
				msg := Response{rooms[roomdId].matrix}
				sendMessage(roomdId, msg, msg)
			}
			go func() { rooms[roomdId].timer <- true }()
		}
	}
}

func handleTimeExit(user *User) {
	exit := false
	var i int = 1
	for !exit {
		select {
		case <-user.exit:
			exit = true
		default:
			if i >= timeExit*1000 {
				user.userId = ""
				exit = true
			}
			time.Sleep(time.Millisecond)
			i++
		}
	}
	close(user.exit)
}

func handleTimer(roomdId string) {
	exit := false
	var i int = 1
	for !exit {
		select {
		case v := <-rooms[roomdId].timer:
			if v {
				i = 1
			} else {
				exit = true
			}
		default:
			if i >= timeToPlay*1000 {
				i = 1
				msg := Message{"time expired"}
				if rooms[roomdId].play == 1 {
					sendMessage(roomdId, msg, "")
					rooms[roomdId].play = 2
				} else {
					sendMessage(roomdId, "", msg)
					rooms[roomdId].play = 1
				}
			} else {
				if i%1000 == 0 {
					sendMessage(roomdId, i/1000, i/1000)
				}
				time.Sleep(time.Millisecond)
				i++
			}
		}
	}
}

func handleTimerStart(roomId string) {
	for i := timStart; i >= 0; i-- {
		time.Sleep(time.Second)
		sendMessage(roomId, i, i)
	}
	go handleTimer(roomId)
}
