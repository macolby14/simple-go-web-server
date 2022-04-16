package main

import (
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

func auth(res http.ResponseWriter, req *http.Request) {
	if user, err := gothic.CompleteUserAuth(res, req); err == nil {
		session, _ := store.Get(req, "app-session")
		session.Values["user"] = user
		session.Save(req, res)
		fmt.Fprintf(res, "Auth already complete %v %v", res, user)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func authCallback(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	session, _ := store.Get(req, "app-session")
	session.Values["user"] = user
	session.Save(req, res)
	fmt.Fprintf(res, "User info %v", user)
}

func authLogout(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)
	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func home(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "app-session")
	username, _ := session.Values["user"]

	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "Home Page. Username is %v", username)
}

var store *sessions.CookieStore

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Errir loading .env file")
	}

	store = sessions.NewCookieStore([]byte(os.Getenv("APP_SESSION_SECRET")))

	/* This is just used for gothic state */
	os.Setenv("SESSION_SECRET", string(securecookie.GenerateRandomKey(32)))

	goth.UseProviders(google.New(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"), os.Getenv("GOOGLE_OAUTH_SECRET"), "http://localhost:8080/auth/google/callback"))

	router := pat.New()
	router.Get("/auth/{provider}/callback", authCallback)
	router.Get("/auth/{provider}/logout", authLogout)
	router.Get("/auth/{provider}", auth)
	router.Get("/", home)
	http.Handle("/", router)

	fmt.Println("Starting webserver...")
	http.ListenAndServe(":8080", nil)
}
