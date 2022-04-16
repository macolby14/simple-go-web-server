package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/pat"
)

func home(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Home Page")
}


func main() {
	router := pat.New()
	router.Get("/", home)
	http.Handle("/", router)

	fmt.Println("Starting webserver...")
	http.ListenAndServe(":8080", nil)
}