package room

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	Message string `json:"message"`
}

type Request struct {
	Play [2]int	`json:"play"`
}

type Response struct {
	Matrix [6][7]int `json:"matrix"`
}

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
}