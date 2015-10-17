package main

import (
	"log"
	"net/http"

	"github.com/bigrooms/vision/tunnel"
	"github.com/gorilla/websocket"
)

var messages chan []byte

func main() {
	var (
		host = "localhost"
		port = "8080"
	)

	messages = make(chan []byte)

	go tunnel.NewRPCServer(messages, host, port)

	http.HandleFunc("/", dispatchHandler)

	log.Println(http.ListenAndServe("localhost:6060", nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func dispatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("couldn't upgrade:", err)
		return
	}

	defer c.Close()

	for {
		err = c.WriteMessage(websocket.TextMessage, <-messages)

		if err != nil {
			log.Println("error:", err)
		}
	}
}
