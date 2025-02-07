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
	logMu     sync.Mutex                  // Mutex for thread-safe writing to log file
	logFile   *os.File                    // File to save chat history
)

func main() {
	port := ":8989"

	if len(os.Args) == 2 {
		port = ":" + os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	fmt.Println("Listening on the port " + port)

	// Open log file in append mode
	logFile, err = os.OpenFile("chat_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	file, _ := os.ReadFile("linux.txt")

	go broadcaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		if len(clients) >= 10 {
			conn.Write([]byte(string("chatroom full..") + "\n"))
			return
		}

		conn.Write([]byte(string(file) + "\n"))
		conn.Write([]byte(string("[ENTER YOUR NAME]: ")))

		go handleClient(conn)

		// // Send chat history before asking for name
		// sendChatHistory(conn)
	}
}

// broadcaster listens for messages and sends them to all clients
func broadcaster() {
	for msg := range broadcast {
		// Log the message to the file
		logToFile(msg)

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

	// Send chat history before asking for name
	sendChatHistory(conn, name)
	logToFile(fmt.Sprintf("%s has joined our chat...\n", name))
	// Notify all users about the new connection
	for conn1 := range clients {
		if conn1 != conn {
			conn1.Write([]byte(fmt.Sprintf("%s has joined our chat...\n", name)))
		}
	}
	namee := "[" + name + "]"
	timestamp = "[" + timestamp + "]"

	// Continuously listen for messages from this user
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if msg == "/exit\n" {
			break
		}

		// Broadcast the message with the user's name
		broadcast <- fmt.Sprintf("%s%s: %s", timestamp, namee, msg)
	}
	// Handle user disconnection
	mu.Lock()
	delete(clients, conn)
	mu.Unlock()

	broadcast <- fmt.Sprintf("%s has left our chat\n", name)
}

// logToFile writes messages to the chat log file
func logToFile(msg string) {
	logMu.Lock()
	defer logMu.Unlock()
	_, err := logFile.WriteString(msg)
	if err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

// sendChatHistory reads and sends past messages to a new user
func sendChatHistory(conn net.Conn, name string) {
	logMu.Lock()
	defer logMu.Unlock()

	// Open the log file for reading
	file, err := os.Open("chat_log.txt")
	if err != nil {
		conn.Write([]byte("[No chat history available]\n"))
		return
	}
	defer file.Close()

	// Read and send file content
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text() + "\n"))
	}
}