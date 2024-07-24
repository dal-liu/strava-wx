package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type WebhookEvent struct {
	Aspect_type string
	Event_time  int
	Object_id   int
	Object_type string
	Owner_id    int
}

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

func handlePost(w http.ResponseWriter, r *http.Request) {
	var we WebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&we); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		if we.Object_type != "activity" || we.Aspect_type != "create" || we.Owner_id != athleteId {
			return
		}
		if err := updateActivity(we.Object_id); err != nil {
			log.Println(err)
		}
	}()

	w.WriteHeader(http.StatusOK)
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
