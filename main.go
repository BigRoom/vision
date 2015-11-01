package main

import (
	"net/http"
	"os"

	"github.com/bigroom/vision/models"
	"github.com/bigroom/vision/tunnel"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/koding/kite"
	"github.com/paked/configure"
	"github.com/paked/restrict"
	log "github.com/sirupsen/logrus"
)

var (
	messages chan tunnel.MessageArgs
	clients  map[string][]*conn

	conf = configure.New()

	dbName    = conf.String("db-name", "postgres", "DB_NAME")
	dbUser    = conf.String("db-user", "postgres", "DB_USER")
	dbPass    = conf.String("db-pass", "postgres", "DB_PASS")
	dbService = conf.String("db-service", "jarvis", "DB_SERVICE")
	dbPort    = conf.String("db-port", "5432", "DB_PORT")

	defaultIRCServer = conf.String("default-irc", "chat.freenode.net", "default IRC host")

	httpAddr = conf.String("http-address", "0.0.0.0", "Which address you want http to bind on")
	rpcAddr  = conf.String("rpc-address", "0.0.0.0", "Which address you want rpc to bind on")
	httpPort = conf.String("http-port", "6060", "Which port you want http to bind on")
	rpcPort  = conf.String("rpc-port", "8080", "Which port you want rpc to bind on")

	sentryDSN = conf.String("sentry-dsn", "", "The sentry DSN you want to use")

	crypto = conf.String("crypto", "/crypto/app.rsa", "Your crypto")

	sentry *raven.Client

	pool *kite.Client
)

func init() {
	conf.Use(configure.NewEnvironment())
	conf.Use(configure.NewFlag())

	conf.Parse()

	messages = make(chan tunnel.MessageArgs)
	clients = make(map[string][]*conn)
}

func main() {
	var err error
	sentry, err = raven.NewClient(*sentryDSN, nil)
	if err != nil {
		log.Println("No sentry:", err)
	}

	restrict.ReadCryptoKey(*crypto)

	models.Init(
		*dbUser,
		*dbPass,
		*dbService,
		*dbPort,
		*dbName,
	)

	go tunnel.NewRPCServer(messages, *rpcAddr, *rpcPort)
	go messageLoop()

	r := setupRouting()

	pool = setupZombies()

	log.Println(gracehttp.Serve(
		&http.Server{
			Addr:    *httpAddr + ":" + *httpPort,
			Handler: r,
		},
	))
}

func setupRouting() *mux.Router {
	r := mux.NewRouter().
		PathPrefix("/api").
		Subrouter().
		StrictSlash(true)

	r.HandleFunc("/users", registerHandler).
		Methods("POST")

	r.HandleFunc("/users", loginHandler).
		Methods("GET")

	r.HandleFunc("/users/me", restrict.R(secretHandler)).
		Methods("GET")

	r.HandleFunc("/servers/{host}/{channel}/scrollback", restrict.R(messagesHandler)).
		Methods("GET")

	r.HandleFunc("/servers/default", defaultServerHandler).
		Methods("GET")

	r.HandleFunc("/ws", restrict.R(dispatchHandler))

	return r
}

func setupZombies() *kite.Client {
	k := kite.New("vision", "1.0.0")

	url := "http://" + os.Getenv("ZOMBIES_PORT_3001_TCP_ADDR") + ":3001/kite"

	pool := k.NewClient(url)
	go func() {
		connected, err := pool.DialForever()
		if err != nil {
			log.Warn("Got error connecting to zombies", err)
		}

		<-connected

		log.Println("Connected!")
	}()

	return pool
}
