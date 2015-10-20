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
		log.Printf("Dispatching message '%v' to channel with key: '%v'", m.Content, m.Key())

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

	server := r.FormValue("server")
	if server == "" {
		server = "chat.freenode.net:6667"
	}

	user, err := newConnection(server, fmt.Sprintf("roombot%v", rand.Intn(999)), c)
	if err != nil {
		log.Println("couldnt create connection", err)
		return
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
			user.messages <- a.Message
		}

	}
}

type action struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Channel string `json:"channel"`
}

func newConnection(server, nick string, c *websocket.Conn) (*conn, error) {
	fmt.Println("Connecting to server", server)
	zombie := ircx.Classic(server, nick)

	if err := zombie.Connect(); err != nil {
		return nil, err
	}

	user := &conn{
		c:        c,
		irc:      zombie,
		messages: make(chan string),
		server:   server,
		nick:     nick,
	}

	zombie.HandleFunc(irc.PING, user.pingHandler)
	zombie.HandleFunc(irc.RPL_WELCOME, user.registerHandler)
	zombie.HandleFunc(irc.JOIN, user.messageHandler)

	go zombie.HandleLoop()

	return user, nil
}

type conn struct {
	c        *websocket.Conn
	irc      *ircx.Bot
	messages chan string
	nick     string
	server   string
}

func (c *conn) changeNick(name string) {
	c.nick = name

	c.irc.Sender.Send(&irc.Message{
		Command: irc.NICK,
		Params:  []string{c.nick},
	})
}

func (c *conn) messageHandler(s ircx.Sender, m *irc.Message) {
	go func() {
		for {
			log.Println("Waiting for message")
			msg := <-c.messages
			log.Printf("Got message '%v'", msg)
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
