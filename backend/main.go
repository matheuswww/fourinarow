package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	user_controller "github.com/matheuswww/fourinarow/controller/user"
	"github.com/matheuswww/fourinarow/mysql"
	"github.com/matheuswww/fourinarow/room"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error trying load .env")
	}
	mysql.NewMysql()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, 
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
}))
	r.GET("/room", room.RoomConn)
	r.GET("/rooms", room.GetRooms)
	r.POST("/signin", user_controller.Signin)
	r.POST("/signup", user_controller.Signup)
	r.Run(":8080")
}