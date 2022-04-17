package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type User struct {
	Provider  string
	Email     string
	Name      string
	AvatarURL string
}

func createSession(user goth.User, res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "app-session")
	if err != nil {
		log.Printf("[ERROR] getting a session: %v\n", err)
		return
	}
	session.Values["user"] = User{Provider: user.Provider, Email: user.Email, Name: user.Name, AvatarURL: user.AvatarURL}
	fmt.Println(session.Values["user"])
	if session.Save(req, res); err != nil {
		fmt.Fprintln(res, "Could not save session", err)
	}
	// res.WriteHeader(http.StatusOK)
	// fmt.Fprintln(res, "Logged in successfully")
}

func auth(res http.ResponseWriter, req *http.Request) {
	if user, err := gothic.CompleteUserAuth(res, req); err == nil {
		createSession(user, res, req)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func authCallback(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(res, err)
		return
	}
	createSession(user, res, req)
}

func authLogout(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)
	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func authHealth(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "app-session")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(map[string]error{"error": err})

	}
	user, _ := session.Values["user"]
	fmt.Println(user)
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]interface{}{"ok": true, "user": user})
}

func health(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]bool{"ok": true})
}

var store *sessions.CookieStore

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	store = sessions.NewCookieStore([]byte(os.Getenv("APP_SESSION_SECRET")))

	/* This is just used for gothic state */
	os.Setenv("SESSION_SECRET", string(securecookie.GenerateRandomKey(32)))

	goth.UseProviders(google.New(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"), os.Getenv("GOOGLE_OAUTH_SECRET"), "http://localhost:8080/auth/google/callback"))

	router := pat.New()
	router.Get("/auth/health", authHealth)
	router.Get("/auth/{provider}/callback", authCallback)
	router.Get("/auth/{provider}/logout", authLogout)
	router.Get("/auth/{provider}", auth)
	router.Get("/api/health", health)
	http.Handle("/", router)

	fmt.Println("Starting webserver...")
	http.ListenAndServe(":8080", nil)
}
