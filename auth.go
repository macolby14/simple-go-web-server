package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var store *sessions.CookieStore

func authInit() {
	gob.Register(User{})

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	goth.UseProviders(google.New(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"), os.Getenv("GOOGLE_OAUTH_SECRET"), "http://localhost:8080/api/auth/google/callback"))
	gothic.Store = store
}

func createSession(user goth.User, res http.ResponseWriter, req *http.Request) {

	appUser := getOrCreateUser(user)
	log.Printf("Temp log. User. %v\n", appUser)

	session, err := store.Get(req, "app-session")
	if err != nil {
		log.Printf("[ERROR] getting a session: %v\n", err)
		return
	}
	session.Values["user"] = User{Email: user.Email, AvatarURL: user.AvatarURL}
	if err = session.Save(req, res); err != nil {
		fmt.Fprintln(res, "Could not save session", err)
	}
	res.Header().Set("Location", "/api/auth/user")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func auth(res http.ResponseWriter, req *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		createSession(gothUser, res, req)
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

func authUser(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "app-session")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(map[string]error{"error": err})

	}
	user, _ := session.Values["user"]
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]interface{}{"ok": true, "user": user})
}
