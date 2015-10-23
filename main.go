package main

import (
	"log"
	"net/http"

	"github.com/bigroom/vision/models"
	"github.com/bigroom/vision/tunnel"
	"github.com/bigroom/zombies"
	"github.com/gorilla/mux"
	"github.com/paked/configure"
	"github.com/paked/restrict"
)

var (
	messages chan tunnel.MessageArgs
	clients  map[string][]*zombies.Zombie

	conf = configure.New()

	dbName    = conf.String("db-name", "postgres", "DB_NAME")
	dbUser    = conf.String("db-user", "postgres", "DB_USER")
	dbPass    = conf.String("db-pass", "postgres", "DB_PASS")
	dbService = conf.String("db-service", "jarvis", "DB_SERVICE")
	dbPort    = conf.String("db-port", "5432", "DB_PORT")

	crypto = conf.String("crypto", "/crypto/app.rsa", "Your crypto")
)

func main() {
	conf.Use(configure.NewEnvironment())
	conf.Use(configure.NewFlag())

	conf.Parse()

	restrict.ReadCryptoKey(*crypto)

	models.Init(
		*dbUser,
		*dbPass,
		*dbService,
		*dbPort,
		*dbName,
	)

	var (
		host = "0.0.0.0"
		port = "8080"
	)

	messages = make(chan tunnel.MessageArgs)
	clients = make(map[string][]*zombies.Zombie)

	go tunnel.NewRPCServer(messages, host, port)
	go messageLoop()

	r := mux.NewRouter()

	r.HandleFunc("/ws", dispatchHandler)

	r.HandleFunc("/users", registerHandler).
		Methods("POST")

	http.Handle("/", r)

	log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
}
