package jwt_service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Token struct {
	Token string `json:"token"`
}

func GetToken(c *gin.Context) {
	userId := uuid.NewString()
	user := User{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			Subject: "name",
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(ExpToken).Unix(),
		},
	}
	token,_ := NewAccessToken(user)
	c.JSON( http.StatusOK, struct{ 
		Token  string `json:"token"` 
		UserId string `json:"user_id"`	
	}{ token,  userId})
}