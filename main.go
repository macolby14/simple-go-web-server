package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func auth(res http.ResponseWriter, req *http.Request) {
	if user, err := gothic.CompleteUserAuth(res, req); err == nil {
		fmt.Fprintf(res, "Auth a success %v %v", res, user)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func home(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "Home Page")
}

func main() {
	goth.UseProviders(google.New("clientKey", "secret", "http://localhost:8080/auth/google/callback"))

	router := pat.New()
	router.Get("/auth/{provider}", auth)
	router.Get("/", home)
	http.Handle("/", router)

	fmt.Println("Starting webserver...")
	http.ListenAndServe(":8080", nil)
}
