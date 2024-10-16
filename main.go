package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
)

const port = ":6379"

func main() {
	// Start a TCP listener on port 6379
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer listener.Close()

	log.Printf("Server started, listening on port %s...", port)

	for {
		// Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a reader to read client input
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer

	for {
		// Read data in chunks until we detect \r\n\r\n
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		// Accumulate the data in the buffer
		buffer.Write(data)
		log.Printf("Received from client till now: %s", buffer.String())

		// Check if the buffer contains \r\n\r\n
		if bytes.Contains(buffer.Bytes(), []byte("\n\n")) {
			// Log the input on the terminal
			log.Printf("Received from client: %s", buffer.String())

			// Write the "success" response to the client
			_, err = conn.Write([]byte("success\n"))
			if err != nil {
				log.Printf("Failed to write to connection: %v", err)
				return
			}

			// Clear the buffer for the next message
			buffer.Reset()
		}
	}
}
