package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bigroom/vision/tunnel"
	"github.com/gorilla/websocket"
	"github.com/nickvanw/ircx"
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
		log.Println("Got message")
		log.Printf("Sending to channel with key: %v", m.Key())

		for _, u := range clients[m.Key()] {
			fmt.Println("Writing message")
			err := u.c.WriteJSON(m)
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
		var a action
		err := user.c.ReadJSON(&a)
		if err != nil {
			log.Println("error reading:", err)
			return
		}

		if a.Name == "SET" {
			log.Println("User joined channel", a.Message)
			clients[a.Message] = append(clients[a.Message], user)
		} else if a.Name == "SEND" {
			log.Printf("Sending message '%s' to channel '%s'", a.Message, a.Channel)
		}

	}
}

type conn struct {
	c   *websocket.Conn
	irc *ircx.Bot
}

type action struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Channel string `json:"channel"`
}
