package redis_test

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

    _, err := chatty.Init()

    if err != nil {
        t.Fatal(err.Error())
    }

    _ = chatty.NewUser("god")

    assertUsersCount(t, emptiness, 1)

    heavens, _ := chatty.NewRoom("god")

    if heavens != "1" {
        t.Fatal("Wrong first room id")
    }

    assertUsersCount(t, emptiness, 0)
    assertUsersCount(t, heavens, 1)

    _ = chatty.NewUser("adam")
    _ = chatty.NewUser("eva")

    assertUsersCount(t, emptiness, 2)

    _ = chatty.EnterRoom("adam", heavens)

    assertUsersCount(t, emptiness, 1)
    assertUsersCount(t, heavens, 2)

    _ = chatty.EnterRoom("eva", heavens)

    assertUsersCount(t, emptiness, 0)
    assertUsersCount(t, heavens, 3)

    earth, _ := chatty.NewRoom("adam")

    assertUsersCount(t, earth, 1)
    assertUsersCount(t, emptiness, 0)
    assertUsersCount(t, heavens, 2)

    _ = chatty.LeaveRoom("eva")
    _ = chatty.EnterRoom("eva", earth)

    assertUsersCount(t, earth, 2)
    assertUsersCount(t, emptiness, 0)
    assertUsersCount(t, heavens, 1)
}
