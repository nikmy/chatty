package tests

import (
	chatty "github.com/nikmy/chatty/pkg"
	"strconv"
	"sync"
	"testing"
)

func racyCounts(t *testing.T, nUsers, nIncr int) {
	handleError := func(err error) {
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	handleError(chatty.Init())

	defer func() {
		handleError(chatty.Finalize())
	}()

	users, err := UsersFactory().SpawnUsers(nUsers)
	handleError(err)

	var wg sync.WaitGroup
	wg.Add(nUsers)
	for i := 0; i < nUsers; i++ {
		go func(user *MockUser) {
			for i := 0; i < nIncr; i++ {
				err := user.SendMsg(strconv.Itoa(i))
				handleError(err)
				_, err = user.UpdateHistory()
				handleError(err)
			}
			wg.Done()
		}(&users[i])
	}
	wg.Wait()

	missed := 0
	for _, user := range users {
		h, err := user.UpdateHistory()
		handleError(err)
		missed += len(h)
	}

	if missed > (nUsers * (nUsers - 1) * nIncr / 2) {
		t.FailNow()
	}
}

func TestStressOneUser(t *testing.T) {
	racyCounts(t, 1, 100)
}

func TestStressTwoUsers(t *testing.T) {
	racyCounts(t, 2, 1000)
}

func TestStressManyUsers(t *testing.T) {
	racyCounts(t, 7, 1000)
}
