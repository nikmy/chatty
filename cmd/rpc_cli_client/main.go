package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"runtime"
	"strings"

	chatty "github.com/nikmy/chatty/pkg"
)

func main() {
	CLEAR_FUNC[runtime.GOOS]()
	fmt.Print(HELLO_MESSAGE)
	r := bufio.NewReader(os.Stdin)
	MainCLI(r)
}

var STATE struct {
	User chatty.ClientState
	Addr string
}

func MainCLI(r *bufio.Reader) {
	for {
		fmt.Print("chatty@<offline>$ ")
		cmd, err := r.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if err != nil {
			fmt.Print("\n\n[ERROR] Input error\n\n")
			continue
		}

		if cmd == "" {
			continue
		}

		if cmd == "quit" {
			break
		}

		switch cmd {
		case "conn":
			ConnCLI(r)
		case "help":
			fmt.Print(MAIN_MANUAL)
		default:
			fmt.Print("\nUnknown command. Type `help` for manual.\n\n")
		}
	}
}

func ConnCLI(r *bufio.Reader) {
	var err error

	fmt.Print("\n\tEnter server ip: ")
	STATE.Addr, err = r.ReadString('\n')
	STATE.Addr = strings.TrimSpace(STATE.Addr)

	if err != nil {
		fmt.Print("\t[ERROR] Input error\n\n")
		return
	}

	var username string
	fmt.Print("\tEnter username: ")
	username, err = r.ReadString('\n')
	username = strings.TrimSpace(username)

	if err != nil {
		fmt.Print("\t[ERROR] Input error\n\n")
		return
	}

	remote, err := jsonrpc.Dial("tcp", STATE.Addr+":38120")
	if err != nil {
		fmt.Print("\t[ERROR] Connection error\n\n")
		return
	}

	defer func() {
		_ = remote.Close()
	}()

	if err = remote.Call("Chatty.NewUser", username, &STATE.User); err != nil {
		fmt.Printf("\t[ERROR] RPC Error: %s\n\n", err.Error())
		return
	}

	defer func() {
		_ = remote.Call("Chatty.LeaveRoom", STATE.User, &STATE.User)
	}()

	fmt.Printf("\t[INFO] Successfully connected to %s\n\n", STATE.Addr)

	var cmd string
	for {
		fmt.Printf("chatty@<%s>$ ", STATE.Addr)
		cmd, err = r.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if err != nil {
			fmt.Print("\n[ERROR] Input error\n\n")
			continue
		}

		if cmd == "" {
			continue
		}

		switch cmd {
		case "quit":
			return
		case "join":
			newState := STATE

			fmt.Print("\n\tEnter chat id: ")
			newState.User.RoomId, err = r.ReadString('\n')
			newState.User.RoomId = strings.TrimSpace(newState.User.RoomId)

			if err != nil {
				fmt.Print("\t[ERROR] Input error\n\n")
				return
			}

			err = remote.Call("Chatty.EnterRoom", newState.User, &STATE.User)
			if err != nil {
				fmt.Printf("\t[ERROR] RPC Error: %s\n\n", err.Error())
				continue
			}
			ChatCLI(r, remote)

		case "chat":
			err = remote.Call("Chatty.NewRoom", STATE.User, &STATE.User)
			if err != nil {
				fmt.Printf("\t[ERROR] RPC Error: %s\n\n", err.Error())
				continue
			}
			fmt.Println()
			ChatCLI(r, remote)
		case "help":
			fmt.Print(CONN_MANUAL)
		default:
			fmt.Print("\nUnknown command. Type `help` for manual.\n\n")
		}
	}
}

func ChatCLI(r *bufio.Reader, remote *rpc.Client) {
	fmt.Printf("\t[INFO] You are now in chat with id: %s\n\n", STATE.User.RoomId)
	defer func() {
		_ = remote.Call("Chatty.LeaveRoom", STATE.User, &STATE.User)
	}()

	var cmd string
	var err error

	for {
		fmt.Printf("chatty@<%s>[%s]$ ", STATE.Addr, STATE.User.RoomId)
		cmd, err = r.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if err != nil {
			fmt.Print("\n[ERROR] Input error\n\n")
			continue
		}

		if cmd == "" {
			continue
		}

		switch cmd {
		case "leave":
			return
		case "send":
			SendMessage(r, remote)
		case "update":
			UpdateHistory(remote)
		case "save":
			DumpHistory(remote)
		case "help":
			fmt.Print(CHAT_MANUAL)
		default:
			fmt.Print("\nUnknown command. Type `help` for manual.\n\n")
		}
	}
}
