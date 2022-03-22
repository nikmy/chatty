# Chatty API

## Users-Rooms API

```go
package chatty

// NewUser
//      action: Creates user with desired ID and adds him to waiting room
//      return: Success code
func NewUser(userId string) bool

// NewRoom
//      action: Creates room and adds keyHolder to it
//      return: New room ID and success code
func NewRoom(keyHolderId string) (string, bool)

// EnterRoom
//      action: Moves user from waiting room to desired room
//      return: Success code
func EnterRoom(userId, roomId string) bool

// LeaveRoom
//      action: Moves user from his current room to waiting room 
//      return: Success code
func LeaveRoom(userId string) bool
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
//      action: Picks up all messages are sent after user entered the room
//      return: History and success code
func PickUpHistory(userId string) ([]Message, bool)
```

## Service API

```go
package chatty

// Init
//      action: Launch chatty service, redis and kafka
//      return: Success code and error
func Init() (bool, error)

// Finalize
//      action: Shutdown all services
//      return: Error
func Finalize() error
```

