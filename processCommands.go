package main

import (
	"log"
)

func executeCommand(command string, args []string, aofReplay bool) string {
	log.Println("executing command: ", command, " with args: ", args, " aof replay: ", aofReplay)
	if handler, ok := handlers[command]; !ok {
		return "-ERR unknown command\r\n"
	} else {
		if !aofReplay {
			addToAOFLogs(command, args)
		}
		return handler(args)
	}
}

func processCommand(data []byte) string {
	command, args, err := getCommandAndArgs(data)
	if err != nil {
		return "-ERR unknown command\r\n"
	}
	return executeCommand(command, args, false)
}
