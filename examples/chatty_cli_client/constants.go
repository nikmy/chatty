package main

import (
	"os"
	"os/exec"
)

//goland:noinspection ALL
const (
	HELLO_MESSAGE string = "\n\n" +
		"\n $$$$$$\\  $$\\                  $$\\     $$\\               \n" +
		"$$  __$$\\ $$ |                 $$ |    $$ |              \n" +
		"$$ /  \\__|$$$$$$$\\   $$$$$$\\ $$$$$$\\ $$$$$$\\   $$\\   $$\\ \n" +
		"$$ |      $$  __$$\\  \\____$$\\\\_$$  _|\\_$$  _|  $$ |  $$ |\n" +
		"$$ |      $$ |  $$ | $$$$$$$ | $$ |    $$ |    $$ |  $$ |\n" +
		"$$ |  $$\\ $$ |  $$ |$$  __$$ | $$ |$$\\ $$ |$$\\ $$ |  $$ |\n" +
		"\\$$$$$$  |$$ |  $$ |\\$$$$$$$ | \\$$$$  |\\$$$$  |\\$$$$$$$ |\n" +
		" \\______/ \\__|  \\__| \\_______|  \\____/  \\____/  \\____$$ |\n" +
		"                                               $$\\   $$ |\n" +
		"                                               \\$$$$$$  |\n" +
		"                                                \\______/ \n\n\n"

	BORDER string = "\n--------------------------------------------------------------------------------\n\n"

	MAIN_MANUAL string = BORDER + "Available commands:\n\n" +
		"\thelp\n\t\tPrints manual\n\n" +
		"\tconn\n\t\tConnect server (opens connect dialog)\n\n" +
		"\tquit\n\t\tQuit the app\n" + BORDER

	CONN_MANUAL string = BORDER + "Available commands:\n\n" +
		"\thelp\n\t\tPrints manual\n\n" +
		"\tchat\n\t\tStart new chat (prints chat id and join it)\n\n" +
		"\tjoin\n\t\tJoin an existing chat (opens joining dialog)\n\n" +
		"\tquit\n\t\tGo offline\n" + BORDER

	CHAT_MANUAL string = BORDER + "Available commands:\n\n" +
		"\thelp\n\t\tPrints manual\n" +
		"\tsend\n\t\tSend message (open typing dialog)\n\n" +
		"\tupdate\n\t\tShow unread messages\n\n" +
		"\tsave\n\t\tSave full chat history in .chatty/history.txt\n\n" +
		"\tleave\n\t\tLeave the chat\n" + BORDER
)

//goland:noinspection ALL
var CLEAR_FUNC map[string]func()

func init() {
	CLEAR_FUNC = make(map[string]func())
	CLEAR_FUNC["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
	CLEAR_FUNC["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
	CLEAR_FUNC["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
