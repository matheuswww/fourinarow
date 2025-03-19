package user_repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	rest_err "github.com/matheuswww/fourinarow/restErr"
	"golang.org/x/crypto/bcrypt"
)

func Signin(user_name, password string, db *sql.DB) (string, *rest_err.RestErr) {
  ctx,cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()
  query := "SELECT user_id, password FROM user WHERE user_name = ?"
  var user_id, encryptedPassword string
  err := db.QueryRowContext(ctx, query, user_name).Scan(&user_id, &encryptedPassword)
  if err != nil {
		if err == sql.ErrNoRows {
			return "", rest_err.NewBadRequestError("user not found")
		}
    log.Println("Error trying QueryRowContext: ", err)
    return "", rest_err.NewInternalServerError("server error")
  }
	err = bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		return "", rest_err.NewBadRequestError("invalid password")
	}
  return user_id, nil
}