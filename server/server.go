package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	broadcast = make(chan string)
	mu        sync.Mutex
	logFile   *os.File
)

// StartServer initializes the TCP chat server
func StartServer(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Println("Listening on the port " + port)

	portnum := port[1:]
	logfileName := fmt.Sprintf("chat_log_%s.txt", portnum)

	logFile, err = os.OpenFile(logfileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	go Broadcaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		if len(clients) >= 10 {
			conn.Write([]byte("Chatroom full...\n"))
			conn.Close()
			continue
		}

		go HandleClient(conn, logfileName)
	}
}
