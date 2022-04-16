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


func hello(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request){
	for name, headers := range req.Header {
		for _, h:= range headers {
			fmt.Fprintf(w, "%v : %v\n", name, h)
		}
	}
}

func main() {
	router := pat.New()
	
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	router.Get("/", home)

	http.Handle("/", router)

	fmt.Println("Starting webserver...")
	http.ListenAndServe(":8080", nil)
}