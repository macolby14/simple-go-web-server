package main

import (
	"database/sql"
	"log"

	"github.com/markbates/goth"
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

	// log.Println(db)

	// testInsertQ := `
	//     INSERT INTO user (name, email, avatarUrl, timeCreated)
	//     VALUES ('test', 'test', 'test', 100);
	// 	SELECT * FROM user;
	// `
	// rows, err := db.Query(testInsertQ)
	// if err != nil {
	// 	log.Printf("[ERROR] Error with db query. %v", err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var (
	// 		id          int
	// 		name        string
	// 		email       string
	// 		avatarUrl   string
	// 		timeCreated int
	// 	)
	// 	if err := rows.Scan(&id, &name, &email, &avatarUrl, &timeCreated); err != nil {
	// 		log.Printf("[ERROR] Error scanning row. %v\n", err)
	// 	}
	// 	log.Printf("Results: %v %v %v %v %v\n", id, name, email, avatarUrl, timeCreated)
	// }

}

func getOrCreateUser(gothUser goth.User) *User {
	user, found := getUser(gothUser.Email)
	if found {
		return user
	}
	return nil
}

func getUser(gothEmail string) (*User, bool) {
	// userQ := `
	// 	SELECT name, email, avatarUrl, timeCreated FROM user WHERE email=?;
	// `
	userQ := `
	SELECT * FROM user WHERE email=?;
`
	rows, err := db.Query(userQ, gothEmail)
	if err != nil {
		log.Printf("[ERROR] Error with db query. %v", err)
	}
	defer rows.Close()

	hasNext := rows.Next()

	if !hasNext {
		log.Printf("[INFO] No user account found in getUser for email. %v\n", gothEmail)
		return nil, false
	}

	if rows.Next() {
		log.Fatalf("[ERROR] Multiple users with same emaila ddress. %v\n", gothEmail)
	}

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

	return &User{"provider", email, avatarUrl}, true

}
