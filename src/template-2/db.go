package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func openDB() *sql.DB {

	db, err := sql.Open("mysql",
		"root:myPassw0rd2@tcp(127.0.0.1:3306)/myFile?parseTime=true")
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
