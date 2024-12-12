package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@%s/%s?%s",
		"root", "asd", // user, password
		"",     // address, empty for localhost:3306
		"test", // db
		"",     // options
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
