package main

import (
	"bufio"
	"bytes"
	"log"
	"math/rand"
	"net"
)

func handleConnection(conn net.Conn) {
	userId := rand.Intn(900000) + 100000
	log.Println("Client connected with id: ", userId)
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer
	for {
		// Read the incoming data
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Error reading from client: %v", err)
			return
		}
		// Accumulate data in the buffer
		buffer.Write(data)
		// Check if we have a complete Redis command
		if isCompleteRedisCommand(buffer.Bytes()) {
			log.Println("Whole command is: ", buffer.String())
			// Process the command
			command, args, err := processCommand(buffer.Bytes())
			// Clear the buffer for the next command
			buffer.Reset()
			if err != nil {
				conn.Write([]byte("-ERR unknown command\r\n"))
			}
			log.Println("command is: ", command)
			log.Println("args are: ", args)
			conn.Write([]byte("+OK\r\n"))
		} else {
			log.Println("Incomplete command is: ", buffer.String())
		}
	}
}
