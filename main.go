package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		handlePost(w, r)
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET\n")
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST\n")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
