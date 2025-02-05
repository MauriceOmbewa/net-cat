package main

import (
	"fmt"
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
}
