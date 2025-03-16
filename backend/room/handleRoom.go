package room

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/gorilla/websocket"
	room_request "github.com/matheuswww/fourinarow/room/request"
	room_response "github.com/matheuswww/fourinarow/room/response"
)

var timeToPlay = 5
var timeStart = 3
var timeExit = 10

var width = 6
var height = 7
var rows = 4

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
						msg := room_response.Messages[1001]
						sendMessage(roomdId, "", msg)
						room.user[0].conn = nil
						room.user[0].exit = make(chan bool)
						handleTimeExit(room.user[0])
					}
					if room.user[1].userId == userId {
						msg := room_response.Messages[1001]
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
		var play room_request.Play
		err = json.Unmarshal(message, &play)
		if err != nil {
			log.Println("Error trying Unmarshal")
			continue
		}
		player := 0
		if room.user[0] != nil && room.user[1] != nil && !room.finished {
			if room.user[0].userId == userId {
				if rooms[roomdId].matrix[play.Play[0]][play.Play[1]] == 2 || rooms[roomdId].play != 1 {
					msg := room_response.Message{ Message: room_response.Messages[1010] }
					sendMessage(roomdId, msg, msg)
					continue
				}
				player = 1
				rooms[roomdId].matrix[play.Play[0]][play.Play[1]] = 1
				rooms[roomdId].play = 2
			}
			if room.user[1].userId == userId {
				if rooms[roomdId].matrix[play.Play[0]][play.Play[1]] == 1 || rooms[roomdId].play != 2 {
					msg := room_response.Message{ Message: room_response.Messages[1010] }
					sendMessage(roomdId, msg, msg)
					continue
				}
				player = 2
				rooms[roomdId].matrix[play.Play[0]][play.Play[1]] = 2
				rooms[roomdId].play = 1
			}
			t := true
			res := handleMatrix(room.matrix, play.Play, player)
			if res {
				room.finished = true
				msg := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1002], room.user[player - 1].userId) }
				sendMessage(roomdId, msg, msg)
				t = false
			}
			msg := room_response.Matrix{ Matrix: rooms[roomdId].matrix }
			sendMessage(roomdId, msg, msg)
			go func() { rooms[roomdId].timer <- t }()
		}
	}
}

func handleMatrix(matrix [6][7]int, play [2]int, player int) bool {
	countLH := 0
	countRH := 0
	countUp := 0
	countDown := 0
	countDiagonal1 := 0
	countDiagonal2 := 0

	for i := 1; i < rows; i++ {
		if play[1]+i <= height-1 && matrix[play[0]][play[1]+i] == player {
			countRH++
		}
		if play[1]-i >= 0 && matrix[play[0]][play[1]-i] == player {
			countLH++
		}

		if play[0]+i <= width-1 && matrix[play[0]+i][play[1]] == player {
			countDown++
		}
		if play[0]-i >= 0 && matrix[play[0]-i][play[1]] == player {
			countUp++
		}

    if play[0]-i >= 0 && play[1]-i >= 0 && matrix[play[0]-i][play[1]-i] == player {
			countDiagonal1++
		}
		if play[0]+i < width && play[1]+i < height && matrix[play[0]+i][play[1]+i] == player {
				countDiagonal1++
		}

		if play[0]+i < width && play[1]-i >= 0 && matrix[play[0]+i][play[1]-i] == player {
				countDiagonal2++
		}
		if play[0]-i >= 0 && play[1]+i < height && matrix[play[0]-i][play[1]+i] == player {
				countDiagonal2++
		}
	}

	if (countRH+countLH == rows-1) || (countUp+countDown == rows-1) || (countDiagonal1 == rows-1) || (countDiagonal2 == rows-1) {
		return true
	}
	return false
}