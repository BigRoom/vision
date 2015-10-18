package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bigrooms/vision/tunnel"
	"github.com/gorilla/websocket"
)

var (
	messages chan tunnel.MessageArgs
	clients  map[string][]*conn
)

func main() {
	var (
		host = "localhost"
		port = "8080"
	)

	messages = make(chan tunnel.MessageArgs)
	clients = make(map[string][]*conn)

	go tunnel.NewRPCServer(messages, host, port)
	go messageLoop()

	http.HandleFunc("/", dispatchHandler)

	log.Println(http.ListenAndServe("localhost:6060", nil))
}

func messageLoop() {
	for {
		log.Println("Waiting on message...")
		m := <-messages
		fmt.Println("Got message")

		for _, u := range clients[m.Key()] {
			fmt.Println("Writing message")
			err := u.c.WriteMessage(websocket.TextMessage, []byte(m.String()))
			if err != nil {
				fmt.Println("error (sending message):", err)
			}
		}
	}
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

	user := &conn{
		c: c,
	}

	for {
		_, message, err := user.c.ReadMessage()
		if err != nil {
			log.Println("error reading:", err)
			return
		}

		msg := string(message)
		if msg[:3] == "SET" {
			log.Println("Adding user to channel")
			channel := msg[3:]
			clients[channel] = append(clients[channel], user)
		}

	}
}

type conn struct {
	c        *websocket.Conn
	channels []string
}
