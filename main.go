package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

// Client represents a connected chat client
type Client struct {
	conn     net.Conn
	name     string
	messages chan string
}

// Server represents our chat server
type Server struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
	messages   []string // Store message history
	mutex      sync.Mutex
}

func main() {
	file, err := os.ReadFile("linux.txt")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(file))

	// var name string

	fmt.Print("[ENTER YOUR NAME]:")

	// fmt.Scanln(&name)
	// fmt.Println("welcome", name)

	ln, _ := net.Listen("tcp", ":8989")
	defer ln.Close()

	server := NewServer()

	go server.run()

	// fmt.Printf("Listening on the port :%s\n", port)
	// fmt.Println("Listening on the port :8989")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		s := server
		clientCount := len(s.clients)
		if clientCount >= 10 {
			conn.Write([]byte("Chat is full (max 10 clients). Please try again later.\n"))
			conn.Close()
			continue
		}

		go server.handleClient(conn)
	}
}
