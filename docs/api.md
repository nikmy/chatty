# Chatty API

## Users-Rooms API

```go
package chatty

// NewUser
//      action: Creates user with desired ID and adds him to waiting room
//      return: Error if occur
func NewUser(userId string) error

// NewRoom
//      action: Creates room and adds keyHolder to it
//      return: New room ID and error if occur
func NewRoom(keyHolderId string) (string, error)

// EnterRoom
//      action: Moves user from waiting room to desired room
//      return: Error if occur
func EnterRoom(userId, roomId string) error

// LeaveRoom
//      action: Moves user from his current room to waiting room 
//      return: Error if occur
func LeaveRoom(userId string) error

// UsersCount
//      return: Number of users in the room with given roomId, or (0, error)
func UsersCount(roomId string) (int, error)
```

## Messenger API

```go
package chatty

// Message struct
type Message struct {
    timestamp uint64
    author    string
    text      string
}

// SendMessage
//      action: Sends given message to each of author's roommates
//      return: True if message was created
func SendMessage(msg Message, roomId string) bool

// PickUpHistory
//      action: Picks up all messages were sent after user entered the room
//      return: History and success code
func PickUpHistory(userId string) ([]Message, bool)
```

## Service API

```go
package chatty

// Init
//      action: Launch chatty service, connected with redis and kafka
//      return: Success code and error
func Init(redisURL, kafkaURL string) (bool, error)

// Finalize
//      action: Shutdown all services
//      return: Error
func Finalize() error
```

