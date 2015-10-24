package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	// "github.com/bigroom/vision/models"
	"github.com/bigroom/zombies"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

func messageLoop() {
	for {
		log.Println("Waiting on message...")
		m := <-messages
		log.Printf("Dispatching message '%v' to channel with key: '%v'", m.Content, m.Key())

		for _, u := range clients[m.Key()] {
			fmt.Println("Writing message")
			err := u.WSConn.WriteJSON(m)
			if err != nil {
				fmt.Println("error (sending message):", err)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func dispatchHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	fmt.Println(t.Claims["id"])

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("couldn't upgrade:", err)
		return
	}

	defer c.Close()

	server := r.FormValue("server")
	if server == "" {
		server = "chat.freenode.net:6667"
	}

	user, err := zombies.New(server, fmt.Sprintf("%v%v", "uUsername", rand.Intn(9999)), c)
	if err != nil {
		log.Println("couldnt create connection", err)
		return
	}

	for {
		var a action
		err := user.WSConn.ReadJSON(&a)
		if err != nil {
			log.Println("error reading:", err)
			return
		}

		if a.Name == "SET" {
			log.Println("User joined channel", a.Message)
			clients[a.Message] = append(clients[a.Message], user)
		} else if a.Name == "SEND" {
			log.Printf("Sending message '%s' to channel '%s'", a.Message, a.Channel)
			user.Messages <- a.Message
		} else if a.Name == "NICK" {
			log.Printf("Changing nick to '%v'", a.Message)
			user.SetNick(a.Message)
		}
	}
}

type action struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Channel string `json:"channel"`
}
