package main

import (
	"log"
	"net"
	"time"
)

const port = ":6379"

func main() {
	flags := getFlags()
	if flags.persistence == "snapshot" {
		log.Println("Persistence mode: snapshot")
		if err := loadLastSnapshot(); err != nil {
			log.Fatal("Snapshot load failed: ", err)
		}
		SetInterval(5*time.Second, createSnapshot)
	} else {
		log.Println("Persistence mode: inmemory")
	}
	addBasicHandlers()
	addListHandlers()
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
