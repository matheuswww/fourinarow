package jwt_service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var ExpToken = time.Hour*4

type Token struct {
	Token string `json:"token"`
}

type User struct {
	UserId string
	jwt.StandardClaims
}

func NewAccessToken(claims User) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) *User {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &User{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil || !parsedAccessToken.Valid {
		return nil
	}
	claims, ok := parsedAccessToken.Claims.(*User)
	if !ok {
		return nil
	}
	return claims
}