package chatty

import (
    "errors"
    "github.com/go-redis/redis"
    chatty "github.com/nikmy/chatty/internal"
    "strconv"
)

func NewUser(userId string) error {
    _, err := usersRC.Get(userId).Int()
    if err != redis.Nil {
        return errors.New("NewUser: given ID is taken")
    }

    err = usersRC.Set(userId, "0", 0).Err()
    if err != nil {
        return err
    }

    err = roomsRC.Incr("0").Err()
    if err != nil {
        return err
    }

    return nil
}

func NewRoom(keyHolderId string) (string, error) {
    nextFreeRoom++
    roomId := strconv.FormatUint(nextFreeRoom, 10)
    roomsRC.Set(roomId, 0, 0)
    return roomId, EnterRoom(keyHolderId, roomId)
}

func EnterRoom(userId, roomId string) error {
    currentRoom, err := usersRC.Get(userId).Result()
    if err == redis.Nil {
        return errors.New("EnterRoom: user does not exists")
    }

    err = roomsRC.Get(roomId).Err()
    if err == redis.Nil {
        return errors.New("EnterRoom: room does not exists")
    }

    err = usersRC.Set(userId, roomId, 0).Err()
    if err != nil {
        return err
    }

    err = roomsRC.Incr(roomId).Err()
    if err != nil {
        return err
    }

    cnt, err := roomsRC.Decr(currentRoom).Result()
    if err != nil {
        return err
    }
    if currentRoom != "0" && cnt == 0 {
        return roomsRC.Del(roomId).Err()
    }

    return nil
}

func LeaveRoom(userId string) error {
    roomId, err := usersRC.Get(userId).Result()
    if err == redis.Nil {
        return errors.New("LeaveRoom: user does not exists")
    }

    if roomId != "0" {
        return EnterRoom(userId, "0")
    }

    return usersRC.Del(userId).Err()
}

func UsersCount(roomId string) (int, error) {
    cnt, err := roomsRC.Get(roomId).Int()
    if err != nil {
        return 0, err
    }
    return cnt, nil
}

var usersRC *redis.Client = nil
var roomsRC *redis.Client = nil

var nextFreeRoom uint64 = 0

func initRedis() (bool, error) {
    suc, err := chatty.ConnectRedis(&usersRC, &roomsRC)
    clearDatabases()
    return suc, err
}

func clearDatabases() {
    if usersRC != nil {
        usersRC.FlushAll()
    }
    if roomsRC != nil {
        roomsRC.FlushAll()
    }
}
