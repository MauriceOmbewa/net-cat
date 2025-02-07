# Net-Cat

![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

A TCP-based chat server implementation in Go, recreating NetCat functionality with added features. This project follows a Server-Client architecture, supporting multiple clients, real-time message broadcasting, and chat history.

## Features

- **Multi-client Support** – Handles multiple concurrent clients
- **Real-time Message Broadcasting** – Instant message delivery to all clients
- **Username Management** – Custom usernames with change capability
- **Chat History** – New clients receive previous messages
- **Join/Leave Notifications** – Informs all users of client activity
- **Activity Logging** – Saves chat logs to a file
- **Scalability** – Efficient use of Goroutines and channels for performance

## Getting Started

### Prerequisites

- Go 1.16 or higher
- netcat (`nc`) for client connections

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/MauriceOmbewa/net-cat.git
   cd net-cat
   ```
2. Build the project:
   ```bash
   go build -o TCPChat
   ```

### Running the Server

```bash
# Default port (8989)
./TCPChat

# Custom port
./TCPChat 2525
```

### Connecting as a Client

```bash
# Using netcat on default port
nc localhost 8989

# Using netcat on custom port
nc localhost 2525
```

## Usage

When a client connects, they will be prompted to enter a username and can start chatting immediately.

### Commands
- `/exit` – Disconnect from the chat
- `/change <new_name>` – Change username

## Project Structure

```
├── LICENSE
├── README.md
├── main.go
├── server/
│   ├── server.go       # Server initialization
│   ├── client.go       # Client handling
│   ├── broadcaster.go  # Message broadcasting
│   ├── logger.go       # Chat history logging
│   ├── utils.go        # Utility functions
```

## Testing

Run the test suite:
```bash
go test ./...
```

## Message Format

Messages are structured as follows:
```
[YYYY-MM-DD HH:MM:SS][username]: message
```
System notifications:
```
[YYYY-MM-DD HH:MM:SS] username has joined the chat...
[YYYY-MM-DD HH:MM:SS] username has left the chat...
```

## Logging

The server maintains a chat log file (`chat_log_<port>.txt`) containing:
- Server start/stop events
- Client connections/disconnections
- Chat messages
- Errors encountered

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

- **Maurice Ombewa** - [GitHub](https://github.com/MauriceOmbewa)

## Acknowledgments

- Inspired by NetCat
- Thanks to all contributors and testers

