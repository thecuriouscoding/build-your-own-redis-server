package main

import (
	"log"
)

func executeCommand(command string, args []string) string {
	log.Println("executing command: ", command, " with args: ", args)
	if handler, ok := handlers[command]; !ok {
		return "-ERR unknown command\r\n"
	} else {
		return handler(args)
	}
}

func processCommand(data []byte) string {
	command, args, err := getCommandAndArgs(data)
	if err != nil {
		return "-ERR unknown command\r\n"
	}
	return executeCommand(command, args)
}
