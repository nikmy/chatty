package main

import (
	"context"
	chatty "github.com/nikmy/chatty/pkg"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os/signal"
	"syscall"
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

func (*Chatty) PickUpHistory(user chatty.ClientState, reply *chatty.UserHistory) (err error) {
	reply.User, reply.History, err = chatty.PickUpHistory(user)
	return
}

func (*Chatty) DumpHistory(user chatty.ClientState, reply *chatty.UserHistory) (err error) {
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
		Logger.Println(err.Error())
		return
	}

	t, err := net.ResolveTCPAddr("tcp", ":"+PORT)
	if err != nil {
		Logger.Println(err.Error())
		return
	}

	l, err := net.ListenTCP("tcp", t)
	if err != nil {
		Logger.Println(err.Error())
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	Logger.Printf("Listening at port %s", PORT)

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				continue
			}
			Logger.Printf("New client: %s", c.RemoteAddr())
			go jsonrpc.ServeConn(c)
		}
	}()

	<-ctx.Done()
	Logger.Println("Graceful shutdown")
}
