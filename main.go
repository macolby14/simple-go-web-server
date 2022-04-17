package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
)

func health(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]bool{"ok": true})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	authInit()

	router := pat.New()
	router.Get("/api/auth/user", authUser)
	router.Get("/api/auth/{provider}/callback", authCallback)
	router.Get("/api/auth/{provider}/logout", authLogout)
	router.Get("/api/auth/{provider}", auth)
	router.Get("/api/health", health)
	http.Handle("/", router)

	log.Println("Starting webserver...")
	http.ListenAndServe(":8080", nil)
}
