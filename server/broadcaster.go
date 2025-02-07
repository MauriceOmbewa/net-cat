package server

// Broadcaster listens for messages and sends them to all clients
func Broadcaster() {
	for msg := range broadcast {
		LogToFile(msg)

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
