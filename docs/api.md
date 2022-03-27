# Chatty API

## Rooms API

```go
package chatty

// ClientState struct
type ClientState struct {
    UserId string
    RoomId string
}

// NewUser
//      action: Creates user with desired ID and adds him to waiting room
//      return: Updated user state and error
func NewUser(userId string) (ClientState, error)

// NewRoom
//      action: Creates room and adds keyHolder to it
//      return: Updated user state and error
func NewRoom(keyHolderId string) (ClientState, error)

// EnterRoom
//      action: Moves user from waiting room to desired room
//      return: Updated user state and error
func EnterRoom(userId, roomId string) (ClientState, error)

// LeaveRoom
//      action: Moves user from his current room to waiting room 
//      return: Updated user state and error
func LeaveRoom(userId string) (ClientState, error)

// UsersCount
//      return: Number of users in the room with given roomId, or (0, error)
func UsersCount(roomId string) (int, error)
```

## Message API

```go
package chatty

// Message struct
type Message struct {
    Timestamp uint64
    Author    string
    Text      string
    Room      string
}

// SendMessage
//      action: Sends given message to each of author's roommates
//      return: Error
func SendMessage(content []byte, user ClientState) error

// PickUpHistory
//      action: Picks up all messages were sent after user entered the room
//      return: History and error
func PickUpHistory(user ClientState) ([]Message, error)

// DumpHistory
//      action: Dumps all messages were sent to the room
//      return: Room history and error
func DumpHistory(user ClientState) ([]Message, error)
```

## Service API

```go
package chatty

// Init
//      action: Launch chatty service, connected with redis and kafka
//      return: Error
func Init(redisURL, kafkaURL string) error

// Finalize
//      action: Shutdown all services
//      return: Error
func Finalize() error
```
