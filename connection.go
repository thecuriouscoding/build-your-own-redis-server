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
			response := processCommand(buffer.Bytes())
			log.Println("response to send is: ", response)
			log.Println("is response -ERR unknown command", response == "-ERR unknown command\r\n")
			// Clear the buffer for the next command
			buffer.Reset()
			conn.Write([]byte(response))
			// conn.Write([]byte("-ERR unknown command\r\n"))
		}
		// else {
		// 	log.Println("Incomplete command is: ", buffer.String())
		// }
	}
}
