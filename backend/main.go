package main

import (
	"log"

	"github.com/gin-contrib/cors"
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
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Permite todas as origens
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Métodos permitidos
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // Headers permitidos
		ExposeHeaders:    []string{"Content-Length"}, // Headers expostos
		AllowCredentials: true, // Permite credenciais (cookies, auth headers, etc.)
}))
	r.GET("/room", room.RoomConn)
	r.GET("rooms", room.GetRooms)
	r.GET("/token", jwt_service.GetToken)
	r.Run(":8080")
}