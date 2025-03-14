package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	jwt_service "github.com/matheuswww/fourinarow/jwt"
	"github.com/matheuswww/fourinarow/room"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error trying load .env")
	}
	r := gin.Default()
	r.GET("/room", room.RoomConn)
	r.GET("/token", jwt_service.GetToken)
	r.Run(":8080")
}