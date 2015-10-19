package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/bigroom/vision/tunnel"
	"github.com/gorilla/websocket"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
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

	zombie := ircx.Classic("chat.freenode.net:6667", fmt.Sprintf("roombot%v", rand.Intn(999)))
	if err != nil {
		log.Println("Could not connect to IRC")
	}

	if err := zombie.Connect(); err != nil {
		log.Println("Unable to connect to IRC", err)
	}

	user := &conn{
		c:        c,
		irc:      zombie,
		messages: make(chan string),
	}

	zombie.HandleFunc(irc.PING, user.pingHandler)
	zombie.HandleFunc(irc.RPL_WELCOME, user.registerHandler)
	zombie.HandleFunc(irc.JOIN, user.messageHandler)

	go zombie.HandleLoop()

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
			user.messages <- a.Message
		}

	}
}

type conn struct {
	c        *websocket.Conn
	irc      *ircx.Bot
	messages chan string
}

type action struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Channel string `json:"channel"`
}

func (c *conn) messageHandler(s ircx.Sender, m *irc.Message) {
	go func() {
		for {
			log.Println("Waiting for message")
			msg := <-c.messages
			log.Println("Got message", msg)
			s.Send(&irc.Message{
				Command:  irc.PRIVMSG,
				Params:   []string{"#roomtest"},
				Trailing: msg,
			})
			log.Println("Message sent")
		}
	}()
}

func (c *conn) registerHandler(s ircx.Sender, m *irc.Message) {
	log.Println("Registering")
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{"#roomtest"},
	})
}

func (c *conn) pingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}
