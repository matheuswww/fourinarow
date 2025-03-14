package jwt_service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GetToken(c *gin.Context) {
	user := User{
		UserId: uuid.NewString(),
		StandardClaims: jwt.StandardClaims{
			Subject: "name",
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(ExpToken).Unix(),
		},
	}
	token,_ := NewAccessToken(user)
	c.JSON( http.StatusOK, struct{ Message string `json:"message"` }{ token })
}