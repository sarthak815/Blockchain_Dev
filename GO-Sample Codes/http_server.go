package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", Home)
	fmt.Println("Starting HTTP Server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Failed to start HTTP Server")
	}
}

type Message struct {
	Dialogues []string `json:"dialogues"`
}
type Request struct {
	MyText []string `json:"mytext"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := new(Request)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		fmt.Println("error here")
		// JSON itself is broken
	}
	// Input Validation
	if len(req.MyText) == 0 {
		fmt.Println("HERE")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message := Message{[]string{
		"Hello There!",
		"General Kenobi!",
		"Welcome to the Party!",
		"Thanks for having me!",
	}}
	message.Dialogues = append(message.Dialogues, req.MyText...)
	bytes, err := json.Marshal(message)
	if err != nil {
		// return 500 error
	}
	w.Write(bytes)
	w.WriteHeader(http.StatusOK)
}
