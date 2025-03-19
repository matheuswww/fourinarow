package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func NewMysql() {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	db := os.Getenv("MYSQL_DB")
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, db)
	var err error
	Db, err = sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
}