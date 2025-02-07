package server

import (
	"net"
)

// NotifyClients sends a message to all clients except the excluded one
func NotifyClients(excludeConn net.Conn, message string) {
	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		if conn != excludeConn {
			conn.Write([]byte(message))
		}
	}
}
