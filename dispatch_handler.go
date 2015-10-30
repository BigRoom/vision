package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bigroom/vision/models"
	"github.com/bigroom/zombies"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

func messageLoop() {
	for {
		log.Println("Waiting for message...")
		m := <-messages
		log.Printf("Dispatching message '%v' to channel with key: '%v'", m.Content, m.Key())
		for _, client := range clients[m.Key()] {
			err := client.c.WriteJSON(response{
				Name:     "MESSAGE",
				Contents: m,
			})

			if err != nil {
				sentry.CaptureErrorAndWait(err, nil)
				fmt.Println("error sending message:", err)
				break
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func dispatchHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	log.Println("Dispatching")
	u, err := models.FetchUser("id", t.Claims["id"])
	if err != nil {
		log.Println("Couldnt get user")
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		sentry.CaptureErrorAndWait(err, nil)
		return
	}

	defer c.Close()

	server := r.FormValue("server")
	if server == "" {
		server = "chat.freenode.net:6667"
	}

	resp, err := pool.Tell("exists", u.ID)
	if err != nil {
		log.Fatal("error getting exists:", err)
		sentry.CaptureErrorAndWait(err, nil)
		return
	}

	if !resp.MustBool() {
		log.Println("Creating zombie")
		add := zombies.Add{
			ID:     u.ID,
			Nick:   u.Username,
			Server: server,
		}

		_, err := pool.Tell("add", add)
		if err != nil {
			log.Fatal("error creating:", err)
			sentry.CaptureErrorAndWait(err, nil)
			return
		}
	}

	for {
		var a action
		err := c.ReadJSON(&a)
		if err != nil {
			log.Println("Closing connection. Error reading:", err)
			return
		}

		if a.Name == "SET" {
			log.Println("Adding user to chanel", a.Message)

			_, err := pool.Tell("join", zombies.Join{
				ID:      u.ID,
				Channel: a.Message,
			})

			if err != nil {
				log.Println("Closing connection. Error joining chanel:", err)
				sentry.CaptureErrorAndWait(err, nil)
				return
			}

			clients[a.Message] = append(clients[a.Message], &conn{
				c:  c,
				id: u.ID,
			})
		} else if a.Name == "SEND" {
			log.Printf("Sending message '%v' to channel '%v'", a.Message, a.Channel)

			_, err := pool.Tell("send", zombies.Send{
				ID:      u.ID,
				Channel: a.Channel,
				Message: a.Message,
			})

			if err != nil {
				log.Println("Closing connection. Error sending message:", err)
				sentry.CaptureErrorAndWait(err, nil)
				return
			}
		} else if a.Name == "CHANNELS" {
			log.Println("Sending channels to user")

			resp, err := pool.Tell("channels", u.ID)
			if err != nil {
				log.Println("Closing connection. Could not connect to kite: ", err)
				sentry.CaptureErrorAndWait(err, nil)
				return
			}

			var channels zombies.Channels
			resp.MustUnmarshal(&channels)

			err = c.WriteJSON(response{
				Contents: channels.Channels,
				Name:     "CHANNELS",
			})

			if err != nil {
				log.Println("Coudlnt wirte JSON")
				sentry.CaptureErrorAndWait(err, nil)
				return
			}
		}
	}
}

type action struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Channel string `json:"channel"`
}

type conn struct {
	id int64
	c  *websocket.Conn
}

type response struct {
	Name     string      `json:"name"`
	Contents interface{} `json:"contents"`
}
