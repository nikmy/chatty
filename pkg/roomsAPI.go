package chatty

import (
	detail "github.com/nikmy/chatty/internal"
)

type ClientState = detail.ClientState

func NewUser(userId string) (ClientState, error) {
	if err := detail.WithRedis().NewUser(userId); err != nil {
		return ClientState{}, err
	}
	return detail.WithKafka().NewUser(userId), nil
}

func NewRoom(keyHolder ClientState) (ClientState, error) {
	room, redisError := detail.WithRedis().NewRoom(keyHolder.UserId)
	if redisError != nil {
		return ClientState{}, redisError
	}
	kafkaError := detail.WithKafka().EnterRoom(&keyHolder, room)
	return keyHolder, kafkaError
}

func EnterRoom(user ClientState, roomId string) (ClientState, error) {
	if err := detail.WithRedis().EnterRoom(user.UserId, roomId); err != nil {
		return ClientState{}, err
	}
	err := detail.WithKafka().EnterRoom(&user, roomId)
	return user, err
}

func LeaveRoom(user ClientState) (ClientState, error) {
	room := user.RoomId
	if err := detail.WithRedis().LeaveRoom(user.UserId); err != nil {
		return user, err
	}
	detail.WithKafka().LeaveRoom(&user)

	if cnt, _ := UsersCount(room); cnt == 0 {
		return user, CloseRoom(room)
	}

	return user, nil
}

func UsersCount(roomId string) (int, error) {
	count, redisError := detail.WithRedis().UsersCount(roomId)
	if redisError != nil {
		return count, redisError
	}
	return count, nil
}

func CloseRoom(roomId string) error {
	if err := detail.WithRedis().CloseRoom(roomId); err != nil {
		return err
	}
	return detail.WithKafka().CloseRoom(roomId)
}
