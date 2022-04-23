package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type User struct {
	Provider  string
	Email     string
	AvatarURL string
}

var db *sql.DB

func dbInit() {
	dsnUri := "/home/mark/simple-go-web-server/db/main.db"
	var err error
	db, err = sql.Open("sqlite", dsnUri)
	if err != nil {
		log.Fatalf("[ERROR] connection to db failed. %v\n", err)
	}

	log.Println(db)

	testInsertQ := `
	    INSERT INTO user (name, email, avatarUrl, timeCreated)
	    VALUES ('test', 'test', 'test', 100);
		SELECT * FROM user;
	`
	rows, err := db.Query(testInsertQ)
	if err != nil {
		log.Printf("[ERROR] Error with db query. %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          int
			name        string
			email       string
			avatarUrl   string
			timeCreated int
		)
		if err := rows.Scan(&id, &name, &email, &avatarUrl, &timeCreated); err != nil {
			log.Printf("[ERROR] Error scanning row. %v\n", err)
		}
		log.Printf("Results: %v %v %v %v %v\n", id, name, email, avatarUrl, timeCreated)
	}

}

// func getOrCreateUser(user goth.User) User {

// }
