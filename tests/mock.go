package tests

import (
	chatty "github.com/nikmy/chatty/pkg"
	"log"
	"strconv"
	"time"
)

type MockUser struct {
	state chatty.ClientState
}

func (u *MockUser) SendMsg(text string) error {
	return chatty.SendMessage([]byte(text), u.state)
}

func (u *MockUser) UpdateHistory() ([]chatty.Message, error) {
	s, h, e := chatty.PickUpHistory(u.state)
	u.state = s
	return h, e
}

func (u *MockUser) SpamDogs(duration time.Duration) {
	go func() {
		for {
			select {
			case _ = <-time.After(duration):
				return
			default:
				if err := u.SendMsg("doggy dog"); err != nil {
					log.Println(err.Error())
					return
				}
			}
		}
	}()
}

type usersFactory struct{}

var mockFactory *usersFactory = nil

func UsersFactory() *usersFactory {
	if mockFactory == nil {
		mockFactory = new(usersFactory)
	}
	return mockFactory
}

func (*usersFactory) SpawnUsers(offset, n int) ([]MockUser, error) {
	first, err := chatty.NewUser(strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	first, err = chatty.NewRoom(first)
	if err != nil {
		return nil, err
	}
	users := []MockUser{{first}}
	room := first.RoomId

	for i := 1; i < n; i++ {
		u, err := chatty.NewUser(strconv.Itoa(i + offset))
		if err != nil {
			return nil, err
		}

		u, err = chatty.EnterRoom(u, room)
		if err != nil {
			return nil, err
		}

		users = append(users, MockUser{u})
	}
	return users, nil
}
