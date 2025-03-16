package room

import (
	"github.com/gorilla/websocket"
)

type User struct {
	conn *websocket.Conn
	userId string
	exit chan bool
}

type Room struct {
	user [2]*User
	play int
	matrix [6][7]int
	timer chan bool
	finished bool
}