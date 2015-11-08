package main

import (
	"net/http"

	"github.com/bigroom/vision/models"
	"github.com/bigroom/zombies"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func messageLoop() {
	for {
		log.Info("Waiting for message...")

		m := <-messages

		log.WithFields(log.Fields{
			"message":     m.Content,
			"channel_key": m.Key(),
		}).Info("Dispatching message to user")

		for _, client := range clients[m.Key()] {
			err := client.c.WriteJSON(response{
				Name:     "MESSAGE",
				Contents: m,
			})

			if err != nil {
				sentry.CaptureError(err, nil)

				log.WithFields(log.Fields{
					"error": err,
				}).Info("Couldnt send error")
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func dispatchHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	u, err := models.FetchUser("id", t.Claims["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"id": t.Claims["id"],
		}).Info("Could not get user with ID")
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		sentry.CaptureError(err, nil)
		return
	}

	defer c.Close()

	server := r.FormValue("server")
	if server == "" {
		server = "chat.freenode.net:6667"
	}

	resp, err := pool.Tell("exists", u.ID)
	if err != nil {
		sentry.CaptureError(err, nil)
		return
	}

	if !resp.MustBool() {
		log.Info("Creating zombie")
		add := zombies.Add{
			ID:     u.ID,
			Nick:   u.Username,
			Server: server,
		}

		_, err := pool.Tell("add", add)
		if err != nil {
			log.Error("error creating:", err)
			sentry.CaptureError(err, nil)
			return
		}
	}

	for {
		var a action
		err := c.ReadJSON(&a)
		if err != nil {
			log.Warn("Closing connection. Error reading:", err)
			return
		}

		if a.Name == "SET" {
			err = handleJoin(a, u, c)
		} else if a.Name == "SEND" {
			err = handleSend(a, u, c)
		} else if a.Name == "CHANNELS" {
			err = handleChannels(a, u, c)
		}

		if err != nil {
			sentry.CaptureError(err, nil)
		}
	}
}

func handleJoin(a action, u models.User, c *websocket.Conn) error {
	log.WithFields(log.Fields{
		"channel_key": a.Message,
	}).Info("Adding user to chanel")

	_, err := pool.Tell("join", zombies.Join{
		ID:      u.ID,
		Channel: a.Message,
	})

	if err != nil {
		log.Warn("Closing connection. Error joining chanel:", err)
		return err
	}

	clients[a.Message] = append(clients[a.Message], &conn{
		c:  c,
		id: u.ID,
	})

	return nil
}

func handleSend(a action, u models.User, c *websocket.Conn) error {
	log.WithFields(log.Fields{
		"message":     a.Message,
		"channel_key": a.Channel,
	}).Info("Going to send message to IRC")

	_, err := pool.Tell("send", zombies.Send{
		ID:      u.ID,
		Channel: a.Channel,
		Message: a.Message,
	})

	if err != nil {
		log.Info("Closing connection. Error sending message:", err)
		return err
	}

	log.WithFields(log.Fields{
		"channel_key": a.Channel,
		"message":     a.Message,
	}).Info("Sent message")

	return nil
}

func handleChannels(a action, u models.User, c *websocket.Conn) error {
	log.Info("Sending channels to user")
	resp, err := pool.Tell("channels", u.ID)
	if err != nil {
		log.Warn("Closing connection. Could not connect to kite: ", err)
		return err
	}

	var channels zombies.Channels
	resp.MustUnmarshal(&channels)

	for _, channel := range channels.Channels {
		key := a.Message + "/" + channel
		log.Info("Joining channel", key)
		clients[key] = append(clients[key], &conn{
			c:  c,
			id: u.ID,
		})
	}

	err = c.WriteJSON(response{
		Contents: channels.Channels,
		Name:     "CHANNELS",
	})

	if err != nil {
		log.Warn("Could not write JSON", err)
		return err
	}

	return nil
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
