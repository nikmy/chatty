package main

import (
	"bufio"
	"fmt"
	chatty "github.com/nikmy/chatty/pkg"
	"net/rpc"
	"os"
)

func SendMessage(r *bufio.Reader, remote *rpc.Client) {
	fmt.Print("\nWrite one-line message, press Enter to send:\n\n> ")
	msg, err := r.ReadBytes('\n')
	if err != nil {
		fmt.Print("\n[ERROR] Input error\n\n")
		return
	}

	if len(msg) > 4096 {
		fmt.Print("\n[ERROR] Message is too long\n\n")
		return
	}

	var sent bool
	err = remote.Call("Chatty.SendMessage", chatty.Message{
		Content: msg,
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

func UpdateHistory(remote *rpc.Client) {
	var reply chatty.UserHistory
	err := remote.Call("Chatty.PickUpHistory", STATE.User, &reply)
	if err != nil {
		fmt.Printf("\n[ERROR] RPC Error: %s\n\n", err.Error())
		return
	}
	STATE.User = reply.User
	saveHistory(reply.History, HISTORY_FILENAME, false)
}

func DumpHistory(remote *rpc.Client) {
	var reply chatty.UserHistory
	err := remote.Call("Chatty.DumpHistory", STATE.User, &reply)
	if err != nil {
		fmt.Printf("\n[ERROR] RPC Error: %s\n\n", err.Error())
		return
	}
	STATE.User = reply.User
	saveHistory(reply.History, HISTORY_FILENAME, true)
}

func header(msg *chatty.Message) []byte {
	y, m, d := msg.Timestamp.Date()
	h, min, s := msg.Timestamp.Clock()
	timeStr := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", y, m, d, h, min, s)
	msgLenStr := fmt.Sprintf("%d", len(msg.Content))
	return []byte(msg.Author + "\t" + timeStr + "\t" + msgLenStr + "\t")
}

func saveHistory(history []chatty.Message, filename string, rewrite bool) {
	flags := os.O_RDWR | os.O_CREATE
	if !rewrite {
		flags |= os.O_APPEND
	}

	file, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		fmt.Printf("\n[ERROR] Cannot open (create) file %s: %s\n\n", filename, err.Error())
	}
	defer func() {
		_ = file.Close()
	}()

	w := bufio.NewWriter(file)
	for _, msg := range history {
		_, _ = w.Write(header(&msg))
		_, _ = w.Write(msg.Content)
	}
	_ = w.Flush()
}
