# Share Note

A WebSocket-based real-time shared editing tool that supports multiple users editing the same text content simultaneously, with real-time synchronization to all connected clients.

## Features

- 📝 **Real-time synchronization**: Real-time data transmission based on WebSocket, ensuring all users see consistent content
- 📱 **Responsive design**: Automatically adapts to PC and mobile devices, providing a good cross-device experience
- 🔌 **Auto-reconnect**: Automatically attempts to reconnect within 30 seconds after network disconnection
- 📊 **Connection status**: Real-time display of current connection status, keeping users informed about synchronization
- ⚙️ **Port configuration**: Supports customizing service port via command-line parameters

## Technical Implementation

- **Backend**: Go language + Gorilla WebSocket library
- **Frontend**: Native HTML + CSS + JavaScript
- **Communication**: WebSocket protocol
- **Concurrency control**: Uses mutex locks to ensure data consistency

## Quick Start

### Prerequisites

- Go 1.16 or higher
- Git

### Installation and Running

1. **Clone the project**

```bash
git clone https://github.com/yourusername/network-toolbox.git
cd network-toolbox/share-note
```

2. **Install dependencies**

```bash
go get github.com/gorilla/websocket
```

3. **Run the service**

```bash
go run main.go
```

By default, the service will start on port 8080. You can access it via browser at `http://localhost:8080`.

### Custom Port

If you need to use a different port, you can specify it via the `-port` parameter:

```bash
go run main.go -port=3000
```

## Usage

1. Open your browser and visit the service address (e.g., `http://localhost:8080`)
2. Enter content in the text box, and all connected users will see updates in real-time
3. The page displays the current connection status to ensure data synchronization is working
4. You can open the same address in multiple browser tabs or different devices to experience real-time synchronization

## Project Structure

```
share-note/
├── main.go          # Main program file
└── README.md        # Project documentation
```

## Core Functionality Implementation

### WebSocket Connection Management

- Establishes WebSocket connections and handles client messages
- Maintains online client list and implements message broadcasting
- Handles connection disconnection and error cases

### Data Synchronization Mechanism

- When a client sends a message, the server will:
  1. Update server-side content
  2. Broadcast the message to all connected clients
  3. Ensure all clients display the same content

### Responsive Design

- Uses CSS media queries to adapt to different screen sizes
- Provides fixed-height editing area on PC
- Automatically adapts to screen height on mobile devices for better touch experience

## Notes

- This project is for demonstration purposes, data persistence is not implemented, data will be lost after service restart
- User authentication and permission management are not implemented, anyone can edit content
- Recommended for use in local network or trusted network environments

## Extension Suggestions

- Add data persistence storage
- Implement user authentication and permission control
- Add text formatting features
- Add history recording and version control
- Implement multi-document support

## License

MIT License

## Contributing

Welcome to submit Issues and Pull Requests to improve this project!