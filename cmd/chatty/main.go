package main

import (
	chatty "github.com/nikmy/chatty/pkg"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Chatty struct{}

func (*Chatty) NewUser(userId string, reply *chatty.ClientState) (err error) {
	*reply, err = chatty.NewUser(userId)
	return
}

func (*Chatty) NewRoom(keyHolder chatty.ClientState, reply *chatty.ClientState) (err error) {
	*reply, err = chatty.NewRoom(keyHolder)
	return
}

func (*Chatty) EnterRoom(user chatty.ClientState, reply *chatty.ClientState) (err error) {
	*reply, err = chatty.EnterRoom(user, user.RoomId)
	if err != nil {
		reply.RoomId = "0"
	}
	return
}

func (*Chatty) LeaveRoom(user chatty.ClientState, reply *chatty.ClientState) (err error) {
	*reply, err = chatty.LeaveRoom(user)
	return
}

func (*Chatty) SendMessage(msg chatty.Message, reply *bool) (err error) {
	err = chatty.SendMessage(msg)
	*reply = err == nil
	return
}

type UserHistory struct {
	User    chatty.ClientState
	History []chatty.Message
}

func (*Chatty) PickUpHistory(user chatty.ClientState, reply *UserHistory) (err error) {
	reply.User, reply.History, err = chatty.PickUpHistory(user)
	return
}

func (*Chatty) DumpHistory(user chatty.ClientState, reply *UserHistory) (err error) {
	reply.History, err = chatty.DumpHistory(user)
	reply.User = user
	return
}

func main() {
	Setup()
	defer Finalize()

	chattyService := new(Chatty)
	err := rpc.Register(chattyService)
	if err != nil {
		log.Println(err.Error())
		return
	}

	t, err := net.ResolveTCPAddr("tcp", ":"+PORT)
	if err != nil {
		log.Println(err.Error())
		return
	}

	l, err := net.ListenTCP("tcp", t)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		log.Printf("Connected to %s", c.RemoteAddr())
		go jsonrpc.ServeConn(c)
	}
}
