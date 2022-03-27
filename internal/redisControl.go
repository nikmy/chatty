package detail

import (
	"errors"
	"github.com/go-redis/redis"
	"strconv"
)

type redisControl struct {
	users        *redis.Client
	rooms        *redis.Client
	nextFreeRoom uint64
}

var globalRedisControl *redisControl = nil

func WithRedis() *redisControl {
	if globalRedisControl != nil {
		return globalRedisControl
	}
	globalRedisControl = new(redisControl)
	*globalRedisControl = redisControl{
		users:        nil,
		rooms:        nil,
		nextFreeRoom: 0,
	}
	return globalRedisControl
}

func (rc *redisControl) NewUser(userId string) error {
	_, err := rc.users.Get(userId).Int()
	if err != redis.Nil {
		return errors.New("NewUser: given ID is taken")
	}

	err = rc.users.Set(userId, "0", 0).Err()
	if err != nil {
		return err
	}

	err = rc.rooms.Incr("0").Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *redisControl) NewRoom(keyHolderId string) (string, error) {
	rc.nextFreeRoom++
	roomId := strconv.FormatUint(rc.nextFreeRoom, 10)
	rc.rooms.Set(roomId, 0, 0)
	return roomId, rc.EnterRoom(keyHolderId, roomId)
}

func (rc *redisControl) EnterRoom(userId, roomId string) error {
	currentRoom, err := rc.users.Get(userId).Result()
	if err == redis.Nil {
		return errors.New("EnterRoom: user does not exists")
	}

	err = rc.rooms.Get(roomId).Err()
	if err == redis.Nil {
		return errors.New("EnterRoom: room does not exists")
	}

	err = rc.users.Set(userId, roomId, 0).Err()
	if err != nil {
		return err
	}

	err = rc.rooms.Incr(roomId).Err()
	if err != nil {
		return err
	}

	cnt, err := rc.rooms.Decr(currentRoom).Result()
	if err != nil {
		return err
	}
	if currentRoom != "0" && cnt == 0 {
		return rc.rooms.Del(roomId).Err()
	}

	return nil
}

func (rc *redisControl) LeaveRoom(userId string) error {
	roomId, err := rc.users.Get(userId).Result()
	if err == redis.Nil {
		return errors.New("LeaveRoom: user does not exists")
	}

	if roomId != "0" {
		return rc.EnterRoom(userId, "0")
	}

	return rc.users.Del(userId).Err()
}

func (rc *redisControl) UsersCount(roomId string) (int, error) {
	cnt, err := rc.rooms.Get(roomId).Int()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (rc *redisControl) CloseRoom(roomId string) error {
	return rc.rooms.Del(roomId).Err()
}

func (rc *redisControl) Init() error {
	err := rc.connectRedis()
	rc.clearDatabases()
	return err
}

func (rc *redisControl) Finalize() error {
	if err := rc.users.Close(); err != nil {
		return err
	}
	if err := rc.rooms.Close(); err != nil {
		return err
	}
	rc.users = nil
	rc.rooms = nil
	rc.nextFreeRoom = 0
	return nil
}

func (rc *redisControl) clearDatabases() {
	if rc.users != nil {
		rc.users.FlushAll()
	}
	if rc.rooms != nil {
		rc.rooms.FlushAll()
	}
}

func (rc *redisControl) connectRedis() error {
	if rc.users != nil && rc.rooms != nil {
		return errors.New("redis connection already established")
	}

	rc.users = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if rc.users == nil {
		return errors.New("cannot connect to users RDB")
	}

	rc.rooms = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	if rc.rooms == nil {
		return errors.New("cannot connect to users RDB")
	}

	return nil
}
