package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func auth(res http.ResponseWriter, req *http.Request) {
	if user, err := gothic.CompleteUserAuth(res, req); err == nil {
		fmt.Fprintf(res, "Auth already complete %v %v", res, user)
	} else {
		fmt.Println("Beginning auth")
		gothic.BeginAuthHandler(res, req)
	}
}

func authCallback(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	fmt.Fprintf(res, "User info %v", user)
}

func logout(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)
	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func home(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "Home Page")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Errir loading .env file")
	}

	os.Setenv("SESSION_SECRET", string(securecookie.GenerateRandomKey(32)))

	goth.UseProviders(google.New(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"), os.Getenv("GOOGLE_OAUTH_SECRET"), "http://localhost:8080/auth/google/callback"))

	router := pat.New()
	router.Get("/auth/{provider}/callback", authCallback)
	router.Get("/auth/{provider}", auth)
	router.Get("/", home)
	http.Handle("/", router)

	fmt.Println("Starting webserver...")
	http.ListenAndServe(":8080", nil)
}
