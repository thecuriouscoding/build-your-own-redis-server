package main

import (
	"bufio"
	"bytes"
	"log"
	"math/rand"
	"net"
)

func handleConnection(conn net.Conn) {
	clientID := rand.Intn(900000) + 100000
	log.Println("Client connected with id: ", clientID)
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
			// Process the command
			response := processCommand(buffer.Bytes())
			// Clear the buffer for the next command
			buffer.Reset()
			conn.Write([]byte(response))
		}
	}
}
