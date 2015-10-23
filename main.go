package main

import (
	"log"
	"net/http"

	"github.com/bigroom/vision/tunnel"
	"github.com/bigroom/zombies"
)

var (
	messages chan tunnel.MessageArgs
	clients  map[string][]*zombies.Zombie
)

func main() {
	var (
		host = "0.0.0.0"
		port = "8080"
	)

	messages = make(chan tunnel.MessageArgs)
	clients = make(map[string][]*zombies.Zombie)

	go tunnel.NewRPCServer(messages, host, port)
	go messageLoop()

	http.HandleFunc("/", dispatchHandler)

	log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
}
