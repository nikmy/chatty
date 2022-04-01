package main

import (
	"flag"
	chatty "github.com/nikmy/chatty/pkg"
	"log"
	"os"
)

// var DETACHED bool
// var QUITE bool

var HELP bool
var PORT string

func init() {
	//flag.BoolVar(&DETACHED, "d", false, "Detached mode: if enabled, logs will be written in .chatty/logs")
	//flag.BoolVar(&QUITE, "q", false, "Quite mode: disable logging for new connections")

	flag.BoolVar(&HELP, "help", false, "Show help message")
	flag.StringVar(&PORT, "port", "38120", "Port")
}

func Setup() {
	flag.Parse()
	if HELP {
		flag.Usage()
		os.Exit(0)
	}

	err := chatty.Init()
	if err != nil {
		log.Fatalf("Cannot initialize chatty:\n%s", err.Error())
	}
}

func Finalize() {
	err := chatty.Finalize()
	if err != nil {
		log.Printf("Cannot finalize chatty:\n%s", err.Error())
	}
}
