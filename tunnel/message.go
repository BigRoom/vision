package tunnel

import (
	"fmt"
	"time"
)

type MessageArgs struct {
	From    string
	Content string
	Time    time.Time
	Channel string
	Host    string
}

func (args MessageArgs) String() string {
	return fmt.Sprintf("[%s] %s (%s) %s {%s}", args.From, args.Content, args.Time, args.Channel, args.Host)
}

type MessageReply struct {
	OK bool
}

type Message struct {
	messages chan []byte
}

func (m *Message) Dispatch(args *MessageArgs, reply *MessageReply) error {
	fmt.Println(args)
	reply.OK = true

	m.messages <- []byte(args.String())

	return nil
}
