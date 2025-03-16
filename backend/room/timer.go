package room

import (
	"time"

	room_response "github.com/matheuswww/fourinarow/room/response"
)

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
			if i > timeToPlay*1000 {
				i = 1
				msg := room_response.Message{Message: "time expired"}
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
	for i := timeStart; i > 0; i-- {
		time.Sleep(time.Second)
		sendMessage(roomId, i, i)
	}
	go handleTimer(roomId)
}
