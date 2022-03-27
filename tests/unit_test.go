package tests

import (
	chatty "github.com/nikmy/chatty/pkg"
	"testing"
)

func assertUsersCount(t *testing.T, roomId string, expectedCnt int) {
	if cnt, _ := chatty.UsersCount(roomId); cnt != expectedCnt {
		t.Fatalf("Wrong users count: %d expected, got %d\n", expectedCnt, cnt)
	}
}

func TestRoomsAPI(t *testing.T) {
	emptiness := "0"

	err := chatty.Init()

	if err != nil {
		t.Fatal(err.Error())
	}

	god, _ := chatty.NewUser("god")

	assertUsersCount(t, emptiness, 1)

	god, _ = chatty.NewRoom(god)
	heavens := god.RoomId

	if heavens != "1" {
		t.Fatal("Wrong first room id")
	}

	assertUsersCount(t, emptiness, 0)
	assertUsersCount(t, heavens, 1)

	adam, _ := chatty.NewUser("adam")
	eva, _ := chatty.NewUser("eva")

	assertUsersCount(t, emptiness, 2)

	adam, _ = chatty.EnterRoom(adam, heavens)

	assertUsersCount(t, emptiness, 1)
	assertUsersCount(t, heavens, 2)

	eva, _ = chatty.EnterRoom(eva, heavens)

	assertUsersCount(t, emptiness, 0)
	assertUsersCount(t, heavens, 3)

	adam, _ = chatty.NewRoom(adam)
	earth := adam.RoomId

	assertUsersCount(t, earth, 1)
	assertUsersCount(t, emptiness, 0)
	assertUsersCount(t, heavens, 2)

	eva, _ = chatty.LeaveRoom(eva)
	eva, _ = chatty.EnterRoom(eva, earth)

	assertUsersCount(t, earth, 2)
	assertUsersCount(t, emptiness, 0)
	assertUsersCount(t, heavens, 1)

	err = chatty.Finalize()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func checkHistory(t *testing.T, actual []chatty.Message, expected []string) {
	actualStr := make([]string, 0)
	for _, m := range actual {
		actualStr = append(actualStr, string(m.Content))
	}
	errWrongHistory := func() {
		t.Fatalf("Wrong history:\n[INFO]\tExpected: %v\n[INFO]\tGot: %v", expected, actualStr)
	}

	if len(actual) != len(expected) {
		errWrongHistory()
	}
	for i, m := range actual {
		if string(m.Content) != expected[i] {
			errWrongHistory()
		}
	}
}

func sendMessages(t *testing.T, sender *chatty.ClientState, mStr []string) {
	for _, m := range mStr {
		err := chatty.SendMessage([]byte(m), *sender)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestMessageAPI(t *testing.T) {
	err := chatty.Init()
	if err != nil {
		t.Fatal(err.Error())
	}

	tester, _ := chatty.NewUser("tester")
	tester, _ = chatty.NewRoom(tester)

	exp := []string{"hi!", "how r u?", "fine!"}

	sendMessages(t, &tester, exp)

	act, err := chatty.DumpHistory(tester)
	checkHistory(t, act, exp)

	err = chatty.Finalize()
	if err != nil {
		t.Fatal(err.Error())
	}
}
