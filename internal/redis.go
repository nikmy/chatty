package chatty

import (
    "errors"
    "github.com/go-redis/redis"
)

func ConnectRedis(usersRC, roomsRC **redis.Client) (bool, error) {
    if *usersRC != nil && *roomsRC != nil {
        return false, errors.New("redis connection already established")
    }

    *usersRC = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    if *usersRC == nil {
        return false, errors.New("cannot connect to users RDB")
    }

    *roomsRC = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       1,
    })

    if *roomsRC == nil {
        return false, errors.New("cannot connect to users RDB")
    }

    return true, nil
}
