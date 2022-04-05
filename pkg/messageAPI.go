package chatty

import (
	detail "github.com/nikmy/chatty/internal"
)

type Message = detail.Message

func SendMessage(msg Message) error {
	return detail.WithKafka().SendMessage(msg)
}

func PickUpHistory(user ClientState) (ClientState, []Message, error) {
	h, e := detail.WithKafka().PickUpHistory(&user)
	return user, h, e
}

func DumpHistory(user ClientState) ([]Message, error) {
	return detail.WithKafka().DumpHistory(&user)
}
