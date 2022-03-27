package chatty

import (
	detail "github.com/nikmy/chatty/internal"
)

type Message = detail.Message

func SendMessage(content []byte, user ClientState) error {
	return detail.WithKafka().SendMessage(content, &user)
}

func PickUpHistory(user ClientState) ([]Message, error) {
	return detail.WithKafka().PickUpHistory(&user)
}

func DumpHistory(user ClientState) ([]Message, error) {
	return detail.WithKafka().DumpHistory(&user)
}
