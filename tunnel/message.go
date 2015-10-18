package tunnel

import (
	"fmt"
	"time"
)

type MessageArgs struct {
	From    string    `json:"from"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
	Channel string    `json:"channel"`
	Host    string    `json:"host"`
}

func (args MessageArgs) Key() string {
	return fmt.Sprintf("%s/%s", args.Host, args.Channel)
}

func (args MessageArgs) String() string {
	return fmt.Sprintf("[%s] %s (%s) %s {%s}", args.From, args.Content, args.Time, args.Channel, args.Host)
}

type MessageReply struct {
	OK bool
}

type Message struct {
	messages chan MessageArgs
}

func (m *Message) Dispatch(args *MessageArgs, reply *MessageReply) error {
	//fmt.Println(args)
	reply.OK = true

	m.messages <- *args
	fmt.Println("Dispatched.")

	return nil
}
