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

// Create a new server instance
func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		messages:   []string{},
	}
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			// Send message history to new client
			for _, msg := range s.messages {
				client.messages <- msg
			}
			s.mutex.Unlock()
			s.broadcast <- fmt.Sprintf("[%s] %s has joined our chat...", time.Now().Format("2006-01-02 15:04:05"), client.name)
			
		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.messages)
			}
			s.mutex.Unlock()
			s.broadcast <- fmt.Sprintf("[%s] %s has left our chat...", time.Now().Format("2006-01-02 15:04:05"), client.name)
			
		case message := <-s.broadcast:
			s.mutex.Lock()
			s.messages = append(s.messages, message)
			for client := range s.clients {
				select {
				case client.messages <- message:
				default:
					close(client.messages)
					delete(s.clients, client)
				}
			}
			s.mutex.Unlock()
		}
	}
}