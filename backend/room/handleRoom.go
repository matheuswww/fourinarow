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

func handleRoom(conn *websocket.Conn, roomId, userId string) {
	room := rooms[roomId]
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				if room.user[0] != nil && room.user[0].conn != nil && room.user[1] != nil && room.user[1].conn != nil {
					if room.user[0].userId == userId {
						msg := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1001], room.user[0].userId) }
						sendMessage(roomId, "", msg)
						room.user[0].conn = nil
						room.user[0].exit = make(chan bool)
						handleTimeExit(room.user[0], room.user[1], roomId)
					}
					if room.user[1].userId == userId {
						msg := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1001], room.user[0].userId) }
						sendMessage(roomId, msg, "")
						room.user[1].conn = nil
						room.user[1].exit = make(chan bool)
						handleTimeExit(room.user[1], room.user[0], roomId)
					}
					break
				}
				if rooms[roomId].start && !rooms[roomId].finished {
					rooms[roomId].timer <- false
				}
				close(rooms[roomId].timer)
				delete(rooms, roomId)
				log.Println(fmt.Sprintf("Room closed: %s", roomId))
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
				player = 1
				if !handlePlay(roomId, play.Play, player, player + 1) {
					continue
				}
			}
			if room.user[1].userId == userId {
				player = 2
				if !handlePlay(roomId, play.Play, player, player - 1) {
					continue
				}
			}
			t := true
			res := handleMatrix(room.matrix, play.Play, player)
			if res {
				room.finished = true
				msg := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1002], room.user[player - 1].userId) }
				sendMessage(roomId, msg, msg)
				log.Println(fmt.Sprintf("Winner: %s, room id: %s", room.user[player - 1].userId, roomId))
				t = false
			}
			msg := room_response.Matrix{ Matrix: rooms[roomId].matrix }
			sendMessage(roomId, msg, msg)
			go func() { rooms[roomId].timer <- t }()
		}
	}
}

func handlePlay(roomId string, play [2]int, player, nextToPlay int) bool {
	if 
	(rooms[roomId].matrix[play[0]][play[1]] != 0 || rooms[roomId].play != player) || 
	(play[0] > width - 1 || play[1] > height - 1) || 
	(play[0] != width - 1 && rooms[roomId].matrix[play[0]+1][play[1]] == 0) || 
	(!rooms[roomId].start) {
		msg := room_response.Message{ Message: room_response.Messages[1010] }
		sendMessage(roomId, msg, msg)
		return false
	}
	rooms[roomId].matrix[play[0]][play[1]] = player
	rooms[roomId].play = nextToPlay
	return true
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