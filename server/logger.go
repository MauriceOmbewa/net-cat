package server

import (
	"bufio"
	"log"
	"net"
	"os"
)

// LogToFile writes messages to the chat log file
func LogToFile(msg string) {
	if logFile == nil {
		return
	}

	_, err := logFile.WriteString(msg)
	if err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

// SendChatHistory sends chat history to a newly connected client
func SendChatHistory(conn net.Conn, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		conn.Write([]byte("[No chat history available]\n"))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text() + "\n"))
	}
}
