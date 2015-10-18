package tunnel

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func NewRPCServer(dispatch chan MessageArgs, host, port string) {
	msg := &Message{
		dispatch,
	}

	rpc.Register(msg)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatal("listen error: ", err)
	}

	log.Panicln(http.Serve(l, nil))
}
