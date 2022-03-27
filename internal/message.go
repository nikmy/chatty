package detail

import (
	"github.com/segmentio/kafka-go"
	"time"
)

type Message struct {
	Timestamp time.Time
	Content   []byte
	Author    string
	Room      string
}

func newMessage(user *ClientState, content []byte) Message {
	return Message{
		Content: content,
		Author:  user.UserId,
		Room:    user.RoomId,
	}
}

func (m *Message) toKafka() kafka.Message {
	return kafka.Message{
		Key:   []byte(m.Author),
		Value: m.Content,
		Time:  m.Timestamp,
	}
}

func (m *Message) fromKafka(km kafka.Message) {
	*m = Message{
		Timestamp: km.Time,
		Content:   km.Value,
		Author:    string(km.Key),
		Room:      km.Topic,
	}
}
