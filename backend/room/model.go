package room

import "github.com/gorilla/websocket"

type Message struct {
	Message string `json:"message"`
}

type Request struct {
	Play  []int	`json:"play"`
}

type Room struct {
	conn *websocket.Conn
	play bool
	userId string
}