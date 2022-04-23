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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "app-session")
		if err != nil {
			log.Printf("[ERROR] getting a session: %v\n", err)
			return
		}
		_, ok := session.Values["user"]
		if ok {
			next.ServeHTTP(res, req)
		} else {
			res.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(res).Encode(map[string]interface{}{"ok": false, "error": "User session not found by middleware for protected route"})
		}
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbInit()
	authInit()

	router := mux.NewRouter()
	protectedRotuer := router.PathPrefix("/api/protected").Subrouter()
	protectedRotuer.Use(authMiddleware)
	router.HandleFunc("/api/auth/user", authUser)
	router.HandleFunc("/api/auth/{provider}/callback", authCallback)
	router.HandleFunc("/api/auth/{provider}/logout", authLogout)
	router.HandleFunc("/api/auth/{provider}", auth)
	router.HandleFunc("/api/health", health)
	protectedRotuer.HandleFunc("/health", health)
	http.Handle("/", router)

	log.Println("Starting webserver...")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
