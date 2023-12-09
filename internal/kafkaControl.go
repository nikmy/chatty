package detail

import (
	"context"
	"errors"
	"io"

	"github.com/segmentio/kafka-go"
)

const kafkaURL = ":9092"

type ClientState struct {
	Offset int64
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
		Offset: 0,
		UserId: userId,
		RoomId: "0",
	}
}

func (kc *kafkaControl) EnterRoom(user *ClientState, roomId string) error {
	if _, ok := kc.rooms[roomId]; !ok {
		newRoom, err := kafka.DialLeader(kafkaCtx, "tcp", kafkaURL, roomId, 0)
		if err != nil {
			return err
		}
		kc.rooms[roomId] = newRoom
	}

	user.Offset, _ = kc.rooms[roomId].ReadLastOffset()
	user.RoomId = roomId
	return nil
}

func (kc *kafkaControl) LeaveRoom(user *ClientState) {
	user.RoomId = "0"
	user.Offset = 0
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

func (kc *kafkaControl) SendMessage(msg Message) error {
	if msg.Room == "0" {
		return errors.New("SendMessage: cannot send messages in wait room")
	}

	conn, ok := kc.rooms[msg.Room]
	if !ok {
		return errors.New("SendMessage: room does not exist")
	}

	_, err := conn.WriteMessages(msg.toKafka())

	return err
}

func (kc *kafkaControl) PickUpHistory(user *ClientState) ([]Message, error) {
	conn, err := kafka.DialLeader(kafkaCtx, "tcp", kafkaURL, user.RoomId, 0)
	if err != nil {
		return nil, err
	}

	_, err = conn.Seek(user.Offset, kafka.SeekStart)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = conn.Close()
	}()

	last, _ := kc.rooms[user.RoomId].ReadLastOffset()
	history, buff := make([]Message, 0), Message{}

	if last <= user.Offset {
		user.Offset = last
		return []Message{}, nil
	}

	for {
		m, err := conn.ReadMessage(1 << 20)
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

	user.Offset = last
	return history, nil
}

func (kc *kafkaControl) DumpHistory(user *ClientState) ([]Message, error) {
	user.Offset = 0
	return kc.PickUpHistory(user)
}

func (kc *kafkaControl) Init() error {
	if kc.rooms != nil {
		return errors.New("kafka controller has been already initialized")
	}

	waitRoom, err := kafka.DialLeader(kafkaCtx, "tcp", kafkaURL, "0", 0)
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
