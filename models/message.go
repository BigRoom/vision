package models

import "time"

const (
	MessagePageSize = 20
	MinOffset       = 0
)

type Message struct {
	ID      int64  `db:"id"`
	Content string `db:"message"`
	User    string `db:"username"`
	Key     string `db:"channel_key"`
	Time    int64  `db:"time"`
}

func NewMessage(content, user, key string) (Message, error) {
	var m Message

	err := DB.
		InsertInto("messages").
		Columns("message", "username", "channel_key", "time").
		Values(content, user, key, time.Now().UnixNano()).
		Returning("*").
		QueryStruct(&m)

	if err != nil {
		return m, err
	}

	return m, nil
}

func Messages(key string, page int64) ([]Message, error) {
	var ms []Message

	err := DB.
		Select("*").
		From("messages").
		Where("channel_key = $1", key).
		OrderBy("time DESC").
		Limit(uint64(MessagePageSize)).
		Offset(uint64(page * MessagePageSize)).
		QueryStructs(&ms)

	if err != nil {
		return ms, err
	}

	return ms, nil
}
