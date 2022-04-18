package main

import (
	"flag"
	chatty "github.com/nikmy/chatty/pkg"
	"log"
	"os"
)

var HELP bool
var PORT string
var LOGFILE string

var logfile = os.Stdout
var Logger *log.Logger

func init() {
	flag.BoolVar(&HELP, "help", false, "Show help message")
	flag.StringVar(&PORT, "port", "38120", "Port")
	flag.StringVar(&LOGFILE, "logfile", "", "File for logging")
}

func Setup() {
	flag.Parse()
	if HELP {
		flag.Usage()
		os.Exit(0)
	}

	if LOGFILE != "" {
		lf, err := os.OpenFile(LOGFILE, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}
		logfile = lf
	}
	Logger = log.New(logfile, "", log.Ldate|log.Ltime)

	if chatty.Init() != nil {
		Logger.Fatal("Cannot initialize chatty")
	}
}

func Finalize() {
	err := chatty.Finalize()
	if err != nil {
		Logger.Printf("Cannot finalize chatty:\n%s", err.Error())
	}
	_ = logfile.Close()
}
