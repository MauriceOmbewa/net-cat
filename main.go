package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	clients   = make(map[net.Conn]string) // Store active clients
	broadcast = make(chan string)         // Channel for broadcasting messages
	mu        sync.Mutex                  // Mutex for thread-safe access to clients map
)

func main() {
	port := ":8989"

	if len(os.Args) == 2 {
		port = ":" + os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("Too many arguments")
		return
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
	}
	defer ln.Close()
	fmt.Println("Listening on the port " + port)

	file, _ := os.ReadFile("linux.txt")

	go broadcaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		if len(clients) >= 10{
			conn.Write([]byte(string("chatroom full..") + "\n"))
			return
		}

		conn.Write([]byte(string(file) + "\n"))
		conn.Write([]byte(string("[ENTER YOUR NAME]: ")))

		// // Add the new connection
		// mu.Lock()
		// clients[conn] = true
		// mu.Unlock()

		// Notify all users about the new connection
		// broadcast <- "New user joined the chat\n"

		go handleClient(conn)

		// conn.Write([]byte(string("New user joined the chat" + "\n")))
	}
}

// broadcaster listens for messages and sends them to all clients
func broadcaster() {
	for msg := range broadcast {
		mu.Lock()
		for conn := range clients {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				conn.Close()
				delete(clients, conn)
			}
		}
		mu.Unlock()
	}
}

// handleClient removes the client when they disconnect
func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1] // Remove newline character
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	mu.Lock()
	clients[conn] = name
	mu.Unlock()

	// Notify all users about the new connection
	broadcast <- fmt.Sprintf("%s has joined our chat...\n", name)
	
	name = "[" + name + "]"
	timestamp = "[" + timestamp + "]"

	// Continuously listen for messages from this user
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// Broadcast the message with the user's name
		broadcast <- fmt.Sprintf("%s%s: %s", timestamp, name, msg)
	}
}
