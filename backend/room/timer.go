package room

import (
	"fmt"
	"time"

	room_response "github.com/matheuswww/fourinarow/room/response"
)

func handleTimeExit(user *User, user2 *User, roomId string) {
	exit := false
	var i int = 0
	for !exit {
		select {
		case <-user.exit:
			exit = true
		default:
			if i > timeExit*1000 {
				user.userId = ""
				exit = true
				msg := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1002], user2.userId) }
				sendMessage(roomId, msg, msg)
			}
			if i % 1000 == 0 {
				msg := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1011], i/1000)}
				sendMessage(roomId, msg, msg)
			}
			time.Sleep(time.Millisecond)
			i++
		}
	}
	close(user.exit)
}

func handleTimer(roomdId string) {
	exit := false
	var i int = 0
	for !exit {
		select {
		case v := <-rooms[roomdId].timer:
			if v {
				i = 0
			} else {
				exit = true
			}
		default:
			if i > timeToPlay*1000 {
				i = 0
				msg := room_response.Message{ Message: room_response.Messages[1000] }
				if rooms[roomdId].play == 1 {
					sendMessage(roomdId, msg, "")
					rooms[roomdId].play = 2
				} else {
					sendMessage(roomdId, "", msg)
					rooms[roomdId].play = 1
				}
			} else {
				if i%1000 == 0 {
					msg := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1012], i/1000) }
					sendMessage(roomdId, msg, msg)
				}
				time.Sleep(time.Millisecond)
				i++
			}
		}
	}
}

func handleTimerStart(roomId string, callback func()) {
	for i := timeStart; i >= 0; i-- {
		time.Sleep(time.Second)
		msg := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1013], i) }
		sendMessage(roomId, msg, msg)
	}
	if rooms[roomId] != nil {
		rooms[roomId].start = true
		callback()
		go handleTimer(roomId)
	}
}
