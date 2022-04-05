package main

import (
	"fmt"
	chatty "github.com/nikmy/chatty/pkg"
	"net/rpc"
)

func SendMessage(remote *rpc.Client) {
	var msg string
	fmt.Print("\nWrite one-line message, press Enter to send:\n\n> ")
	if _, err := fmt.Scanf("%s", &msg); err != nil {
		fmt.Print("\n[ERROR] Input error\n\n")
		return
	}

	if len(msg) > 4096 {
		fmt.Print("\n[ERROR] Message is too long\n\n")
		return
	}

	var sent bool
	err := remote.Call("Chatty.SendMessage", chatty.Message{
		Content: []byte(msg),
		Author:  STATE.User.UserId,
		Room:    STATE.User.RoomId,
	}, &sent)

	if err != nil {
		fmt.Printf("\n[ERROR] RPC Error: %s\n\n", err.Error())
	}

	if sent {
		fmt.Print("\n[INFO] Message was delivered\n\n")
	} else {
		fmt.Print("\n[ERROR] Failed to send message\n\n")
	}
}
