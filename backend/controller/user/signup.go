package user_controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	jwt_service "github.com/matheuswww/fourinarow/jwt"
	user_request "github.com/matheuswww/fourinarow/model/request/user"
	user_response "github.com/matheuswww/fourinarow/model/response/user"
	"github.com/matheuswww/fourinarow/mysql"
	user_repository "github.com/matheuswww/fourinarow/repository/user"
	rest_err "github.com/matheuswww/fourinarow/restErr"
)

func Signup(c *gin.Context) {
	var loginRequest user_request.Login
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		c.JSON(restErr.Code, restErr)
		return
	}
	user_id, restErr := user_repository.Signup(loginRequest.UserName, loginRequest.Password, mysql.Db)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	token, err := jwt_service.NewAccessToken(jwt_service.User{
		UserId: user_id,
		StandardClaims: jwt.StandardClaims{
			Subject: loginRequest.UserName,
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(jwt_service.ExpToken).Unix(),
	}})
	if err != nil {
		restErr := rest_err.NewBadRequestError("server error")
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusOK, user_response.Token{
		Token: token,
		UserId: user_id,
	})
}