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
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Error reading from client: %v", err)
			return
		}
		buffer.Write(data)
		if isCompleteRedisCommand(buffer.Bytes()) {
			response := processCommand(buffer.Bytes())
			buffer.Reset()
			conn.Write([]byte(response))
		}
	}
}
