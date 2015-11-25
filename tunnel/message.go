package tunnel

import (
	"fmt"
	"time"

	"github.com/bigroom/vision/models"
)

type MessageArgs struct {
	ID      int64     `json:"id"`
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

func (args MessageArgs) Message() models.Message {
	return models.Message{
		ID:      args.ID,
		Content: args.Content,
		User:    args.From,
		Key:     args.Key(),
		Time:    args.Time.UnixNano(),
	}
}

type MessageReply struct {
	OK bool
}

type Message struct {
	messages chan MessageArgs
}

func (m *Message) Dispatch(args *MessageArgs, reply *MessageReply) error {
	reply.OK = true

	m.messages <- *args
	return nil
}
