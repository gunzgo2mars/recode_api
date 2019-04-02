package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func IgnitionStart() *sql.DB {

	db, dbErr := sql.Open("mysql", "dev:123456@tcp(127.0.0.1:3306)/recode")

	if dbErr != nil {
		return nil
	}

	return db

}
