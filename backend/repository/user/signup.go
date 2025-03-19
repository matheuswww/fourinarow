package user_repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	rest_err "github.com/matheuswww/fourinarow/restErr"
	"golang.org/x/crypto/bcrypt"
)

func Signup(user_name, password string, sql *sql.DB) (string, *rest_err.RestErr) {
  ctx,cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()
  id := uuid.NewString()
  query := "SELECT COUNT(*) FROM user WHERE user_name = ?"
  var count int
  err := sql.QueryRowContext(ctx, query, user_name).Scan(&count)
  if err != nil {
    log.Println("Error trying QueryRowContext: ", err)
    return "", rest_err.NewInternalServerError("server error")
  }
  if count > 0 {
    return "", rest_err.NewBadRequestError("user already exists")
  }
  encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
  if err != nil {
    log.Println("Error trying GenerateFromPassword: ", err)
    return "", rest_err.NewInternalServerError("server error")
  };
  query = "INSERT INTO user (user_id, user_name, password) VALUES (?, ?, ?)"
  _,err = sql.ExecContext(ctx, query, id, user_name, encryptedPassword)
  if err != nil {
    log.Println("Error trying ExecContext: ", err)
    return "", rest_err.NewInternalServerError("server error")
  }
  return id, nil
}