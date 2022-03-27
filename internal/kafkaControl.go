package detail

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"io"
)

type ClientState struct {
	offset int64
	UserId string
	RoomId string
}

type kafkaControl struct {
	rooms map[string]*kafka.Conn
}

var globalKafkaControl *kafkaControl = nil
var kafkaCtx = context.Background()

func WithKafka() *kafkaControl {
	if globalKafkaControl != nil {
		return globalKafkaControl
	}
	globalKafkaControl = new(kafkaControl)
	*globalKafkaControl = kafkaControl{
		rooms: nil,
	}
	return globalKafkaControl
}

func (kc *kafkaControl) NewUser(userId string) ClientState {
	return ClientState{
		offset: 0,
		UserId: userId,
		RoomId: "0",
	}
}

func (kc *kafkaControl) EnterRoom(user *ClientState, roomId string) error {
	if _, ok := kc.rooms[roomId]; !ok {
		newRoom, err := kafka.DialLeader(kafkaCtx, "tcp", ":9092", roomId, 0)
		if err != nil {
			return err
		}
		kc.rooms[roomId] = newRoom
	}

	user.offset = 0
	user.RoomId = roomId
	return nil
}

func (kc *kafkaControl) LeaveRoom(user *ClientState) {
	user.RoomId = "0"
	user.offset = 0
}

func (kc *kafkaControl) CloseRoom(roomId string) error {
	conn, ok := kc.rooms[roomId]
	if ok {
		if err := conn.DeleteTopics(roomId); err != nil {
			return err
		}
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return errors.New("CloseRoom: kafka: no connection with room")
}

func (kc *kafkaControl) SendMessage(content []byte, user *ClientState) error {
	if user.RoomId == "0" {
		return errors.New("SendMessage: cannot send messages in wait room")
	}

	conn, ok := kc.rooms[user.RoomId]
	if !ok {
		return errors.New("SendMessage: room does not exist")
	}

	m := newMessage(user, content)
	_, err := conn.WriteMessages(m.toKafka())

	return err
}

func (kc *kafkaControl) PickUpHistory(user *ClientState) ([]Message, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{":9092"},
		Topic:     user.RoomId,
		Partition: 0,
		MinBytes:  1 << 10,
		MaxBytes:  1 << 20,
	})
	_ = r.SetOffset(user.offset)

	defer func() {
		_ = r.Close()
	}()

	last, _ := kc.rooms[user.RoomId].ReadLastOffset()
	history, buff := make([]Message, 0), Message{}

	if last <= user.offset {
		user.offset = last
		return []Message{}, nil
	}

	for {
		m, err := kc.rooms[user.RoomId].ReadMessage(1 << 20)
		if err == io.EOF {
			break
		}
		if err != nil {
			return history, err
		}

		buff.fromKafka(m)
		history = append(history, buff)

		if m.Offset == last-1 {
			break
		}
	}

	user.offset = last
	return history, nil
}

func (kc *kafkaControl) DumpHistory(user *ClientState) ([]Message, error) {
	user.offset = 0
	return kc.PickUpHistory(user)
}

func (kc *kafkaControl) Init() error {
	if kc.rooms != nil {
		return errors.New("kafka controller has been already initialized")
	}

	waitRoom, err := kafka.DialLeader(kafkaCtx, "tcp", ":9092", "0", 0)
	if err != nil {
		return err
	}

	kc.rooms = map[string]*kafka.Conn{
		"0": waitRoom,
	}
	return nil
}

func (kc *kafkaControl) Finalize() error {
	if kc.rooms == nil {
		return errors.New("kafka controller has not been initialized")
	}

	for i, room := range kc.rooms {
		if err := room.DeleteTopics(i); err != nil {
			return err
		}
		if err := room.Close(); err != nil {
			return err
		}
	}
	kc.rooms = nil
	return nil
}
