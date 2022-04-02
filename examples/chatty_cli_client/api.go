package main

import (
	"fmt"
	chatty "github.com/nikmy/chatty/pkg"
	"net/rpc"
)

type UserMessage struct {
	User    chatty.ClientState
	Content []byte
}

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

	err := remote.Call("Chatty.SendMessage", UserMessage{
		User:    STATE.User,
		Content: []byte(msg),
	}, &STATE.User)

	if err != nil {
		fmt.Printf("\n[ERROR] RPC Error\n%s\n", err.Error())
	}
}
