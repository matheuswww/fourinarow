package room

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	jwt_service "github.com/matheuswww/fourinarow/jwt"
	room_response "github.com/matheuswww/fourinarow/room/response"
)

var rooms = make(map[string]*Room)

func RoomConn(c *gin.Context) {
	var upgrader = websocket.Upgrader {
		CheckOrigin: func(r *http.Request) bool { 
			return true 
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error trying create connection: ", err)
		return
	}
	defer conn.Close()
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error trying ReadMessage")
		conn.WriteJSON(room_response.Message{ Message: room_response.Messages[1003] })
		conn.Close()
		return
	}
	var token jwt_service.Token
	err = json.Unmarshal(message, &token)
	if err != nil || token.Token == "" {
		log.Println("Error trying Unmarshal")
		conn.WriteJSON(room_response.Message{ Message: room_response.Messages[1003] })
		conn.Close()
		return
	}
	user := jwt_service.ParseAccessToken(token.Token)
	if user == nil {
		log.Println("Invalid token")
		conn.WriteJSON(room_response.Message{ Message: room_response.Messages[1004]})
		conn.Close()
		return
	}
	roomId := c.Request.URL.Query().Get("id")
	if roomId == "" {
		roomId = uuid.NewString()
		err = createRoom(conn, user.UserId, roomId)
	} else if c.Request.URL.Query().Get("rejoin") == "true" {
		err = rejoinRoom(conn, roomId, user.UserId)
	} else {
		err = joinRoom(conn, roomId, user.UserId)
	}
	if err != nil {
		return
	}
	handleRoom(conn, roomId, user.UserId)
}

func createRoom(conn *websocket.Conn, userId, roomId string) error {
	user := &User{
		conn: conn,
		userId: userId,
	}
	timer := make(chan bool)
	room := &Room{
		matrix: [6][7]int{},
		user: [2]*User{user},
		timer: timer,
		finished: false,
	}
	rooms[roomId] = room
	err := conn.WriteJSON(room_response.Message{ Message: fmt.Sprintf("Room_id: %s", roomId) })
	if err != nil {
		log.Println("Error trying WriteJson: ", err)
		conn.Close()
		return err
	}
	log.Println("Room id: ", roomId)
	return nil
}

func joinRoom(conn *websocket.Conn, roomId string, userId string) error {
	r,ok := rooms[roomId]
	if !ok || (r.user[0] != nil && r.user[1] != nil) {
		var msg string
		if !ok {
			msg = room_response.Messages[1007]
		} else {
			msg = "The room is full"
		}
		log.Println(msg)
		conn.WriteJSON(room_response.Message{ Message: msg })
		conn.Close()
		return errors.New(msg)
	}
	if rooms[roomId].user[0].userId == userId {
		msg := room_response.Messages[1005]
		log.Println(msg)
		conn.WriteJSON(room_response.Message{ Message: msg })
		conn.Close()
		return errors.New(msg)
	}
	user := &User{
		conn: conn,
		userId: userId,
		exit: nil,
	}
	r.user[1] = user
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := rd.Intn(2) + 1
	r.play = num
	rooms[roomId] = r
	timeToPlay := r.user[num - 1].userId
	msg := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1008], timeToPlay) }
	sendMessage(roomId, msg, msg)
	go handleTimerStart(roomId)
	log.Println("JoinRoom id: ", roomId)
	return nil
}

func rejoinRoom(conn *websocket.Conn, roomId string, userId string) error {
	r,ok := rooms[roomId]
	if !ok || (r.user[0] == nil || r.user[1] == nil) {
		msg := room_response.Messages[1007]
		log.Println(msg)
		conn.WriteJSON(room_response.Message{ Message: msg })
		conn.Close()
		return errors.New(msg)
	}
	var found = false
	var timeToPlay string
	if r.user[0].userId == userId {
		if r.play == 1 {
			timeToPlay = r.user[0].userId
		}
		r.user[0].conn = conn
		r.user[0].exit <- false
		close(r.user[0].exit)
		found = true
	}
	if r.user[1].userId == userId {
		if r.play == 2 {
			timeToPlay = r.user[1].userId
		}
		r.user[1].conn = conn
		r.user[1].exit <- false
		close(r.user[1].exit)
		found = true
	}
	if !found {
		msg := room_response.Messages[1007]
		log.Println(msg)
		conn.WriteJSON(room_response.Message{ Message: msg })
		conn.Close()
		return errors.New(msg)
	}
	msg := room_response.Matrix{ Matrix: rooms[roomId].matrix }
	sendMessage(roomId, msg, msg)
	msg2 := room_response.Message{ Message: fmt.Sprintf(room_response.Messages[1008], timeToPlay) }
	sendMessage(roomId, msg2, msg2)
	log.Println("Rejoin id: ", roomId)
	return nil
}