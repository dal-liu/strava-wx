package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		handlePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if q.Get("hub.mode") == "subscribe" && q.Get("hub.verify_token") == os.Getenv("VERIFY_TOKEN") {
		w.Write([]byte(q.Get("hub.challenge")))
	}
}

func handlePost(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "POST\n")
}

func main() {
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}

	// err := connectToMongoDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// defer func() {
	// 	if err := disconnectFromMongoDB(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	http.HandleFunc("/webhook", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
