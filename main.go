package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	router.HandleFunc("/api/auth/user", authUser)
	router.HandleFunc("/api/auth/{provider}/callback", authCallback)
	router.HandleFunc("/api/auth/{provider}/logout", authLogout)
	router.HandleFunc("/api/auth/{provider}", auth)
	router.HandleFunc("/api/health", health)
	http.Handle("/", router)

	log.Println("Starting webserver...")
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
