package main

import (
	"database/sql"
	"log"
)

type User struct {
	Provider  string
	Email     string
	AvatarURL string
}

var db *sql.DB

func dbInit() {
	dsnUri := "sqlite:./db/main.db"
	var err error
	db, err = sql.Open("sqlite", dsnUri)
	if err != nil {
		log.Fatalf("[ERROR] connection to db failed. %v\n", err)
	}
}

// func getOrCreateUser(user goth.User) User {

// }
