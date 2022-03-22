# Chatty &mdash; minimalistic chat service

## Requirements:
- Redis for users-rooms cashing
- Apache Kafka for messaging

## Scratch

```go
package scratch

import "fmt"

/*
   package kafka
*/

type Kafka interface {
    NewPartition(topic Topic) Partition
    NewProducer(tp TopicPartition) Producer
    NewConsumer(tp TopicPartition) Consumer
    Message(text string) Message
}

type Message interface{}
type Producer interface{ Produce(msg Message) }
type Consumer interface{ Consume() []Message }

type KafkaClient struct {
    p Producer
    c Consumer
}

func (c KafkaClient) Produce(msg Message) {
    c.p.Produce(msg)
}

func (c KafkaClient) Consume() []Message {
    return c.c.Consume()
}

type TopicPartition struct {
    topic     Topic
    partition Partition
}

type Topic string
type Partition struct {
    topic Topic
    parId uint8
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
   package redis
*/

type Redis interface {
    call(args ...interface{}) interface{}
}

// Next free room ID
var freeRoomId uint64 = 1 // 0 - wait room

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
   package client
*/

/// Globals
///////////////////////////////////////////////////////////////////////////////////////////////////

// Connection with server
var conn interface {
    Request(text string) Client
    Close()
}

// Client state
type Client struct {
    kafka KafkaClient
    name  string
    room  string
}

func newClient(userId, roomId string) Client {
    roomTopic := Topic(roomId)
    userPar := kafka.NewPartition(roomTopic)
    tp := TopicPartition{roomTopic, userPar}
    return Client{
        KafkaClient{
            kafka.NewProducer(tp),
            kafka.NewConsumer(tp),
        },
        userId,
        roomId,
    }
}

// Possible userData state
const (
    Offline = iota
    Online  = iota
    InRoom  = iota
)

// User state struct
var userData struct {
    client Client
    state  uint64
}

///////////////////////////////////////////////////////////////////////////////////////////////////

/// User-room API
///////////////////////////////////////////////////////////////////////////////////////////////////

// Change state to ONLINE and save redis id (offline mode)
func reg(userId string) {
    userData.client = conn.Request(fmt.Sprintf("/reg?userId=%v", userId))
    userData.state = Online
}

// Open new room and enter it (online mode)
func open() {
    userData.client = conn.Request(fmt.Sprintf("/open?keyHolderId=%v", userData.client.name))
    userData.state = InRoom
}

// Enter room (online mode)
func enter(roomId uint64) {
    userData.client = conn.Request(fmt.Sprintf("/enter?roomId=%v&userId=%s", roomId, userData.client.name))
    userData.state = InRoom
}

// Leave room (inRoom mode)
func leave() {
    userData.client = conn.Request(fmt.Sprintf("/leave?roomId=%s&userId=%s", userData.client.room, userData.client.name))
    userData.state = Online
}

// Quit
func quit() {
    switch userData.state {
    case InRoom:
        leave()
        fallthrough
    case Online:
        conn.Request(fmt.Sprintf("/quit?user_id=%v", userData.client.name))
        conn.Close()
        userData.state = Offline
    }
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
   package server
*/

/// Globals
///////////////////////////////////////////////////////////////////////////////////////////////////

/*
   Redis DataBase
       0: user (string) --> room (string)
       1: room (string) --> cnt  (uint64)
*/
var rdb Redis

// Redis command type
type redisCmd interface{}

// Redis command buffer
type redisCB chan redisCmd

// Create redis command
func newRedisCmd(args ...string) redisCmd

/*
   Redis buffered channel associated with connection with usersDB
       To initialize: rdb.call("SELECT 0")
*/
var usersRCB redisCB

/*
   Redis buffered channel associated with connection with usersDB
       To initialize: rdb.call("SELECT 1")
*/
var roomsRCB redisCB

/*
   Buffered channel for results of DECR commands in roomsDB
*/
var roomCnt chan uint64

/*
   Apache Kafka
       roomId --> Topic
       userId --> Partition
*/
var kafka Kafka

///////////////////////////////////////////////////////////////////////////////////////////////////

/// User-room API
///////////////////////////////////////////////////////////////////////////////////////////////////

// /reg?username=...
func regHandler(username string) {
    roomsRCB <- newRedisCmd("SET", username, "0")
}

// /open?keyHolderId=...
func openHandler(keyHolderId string) Client {
    roomId := fmt.Sprintf("%08d", freeRoomId)
    freeRoomId++

    usersRCB <- newRedisCmd("SET", keyHolderId, roomId)
    roomsRCB <- newRedisCmd("SET", roomId, "1")

    return newClient(keyHolderId, roomId)
}

// /enter?userId=...&roomId=...
func enterHandler(userId, roomId string) Client {
    usersRCB <- newRedisCmd("SET", userId, roomId)
    roomsRCB <- newRedisCmd("INCR", roomId)

    return newClient(userId, roomId)
}

// /leave?userId=...&roomId=...
func leaveHandler(roomId, userId string) {
    usersRCB <- newRedisCmd("SET", userId, "0")

    roomsRCB <- newRedisCmd("DECR", roomId)
    if rCnt, _ := <-roomCnt; rCnt == 0 {
        roomsRCB <- newRedisCmd("DEL", roomId)
    }
}

// /quit?userId=...
func quitHandler(userId string) {
    usersRCB <- newRedisCmd("DEL", userId)
}

///////////////////////////////////////////////////////////////////////////////////////////////////

/// Chat API
///////////////////////////////////////////////////////////////////////////////////////////////////

// /update
func updateHandler(user KafkaClient) []Message {
    return user.c.Consume()
}

// /msg
func msgHandler(author KafkaClient, text string) []Message {
    history := author.c.Consume()
    author.p.Produce(kafka.Message(text))
    return history
}

///////////////////////////////////////////////////////////////////////////////////////////////////
```