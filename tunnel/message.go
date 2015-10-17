package tunnel

import (
	"fmt"
	"time"
)

type MessageArgs struct {
	From    string
	Content string
	Time    time.Time
}

func (args MessageArgs) String() string {
	return fmt.Sprintf("[%s] %s (%s)", args.From, args.Content, args.Time)
}

type MessageReply struct {
	OK bool
}

type Message struct{}

func (m *Message) Dispatch(args *MessageArgs, reply *MessageReply) error {
	fmt.Println(args)
	reply.OK = true

	return nil
}
