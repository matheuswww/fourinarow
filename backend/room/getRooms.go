package room

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRooms(c *gin.Context) {
	var roomsId []any
	for k,_ := range rooms {
		roomsId = append(roomsId, struct{ Id string `json:"id"`}{ k })
	}
	c.JSON(http.StatusOK, roomsId)
}