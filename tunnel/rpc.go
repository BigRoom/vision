package tunnel

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func NewRPCServer() {
	msg := &Message{}
	rpc.Register(msg)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", "localhost:6060")
	if err != nil {
		log.Fatal("listen error: ", err)
	}

	log.Panicln(http.Serve(l, nil))
}
