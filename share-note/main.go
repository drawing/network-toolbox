package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	content     string
	clients     = make(map[*websocket.Conn]string) // conn -> clientId
	isLocked    bool
	lockedBy    string
	mutex       sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Message 定义WebSocket消息结构
type Message struct {
	Type     string `json:"type"`
	Content  string `json:"content,omitempty"`
	ClientID string `json:"clientId,omitempty"`
}

// LockStatus 定义锁状态消息
type LockStatus struct {
	Type     string `json:"type"`
	Locked   bool   `json:"locked"`
	LockedBy string `json:"lockedBy,omitempty"`
}

// broadcast 广播消息给所有客户端（可选排除特定客户端）
func broadcast(msg []byte, excludeConn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	
	for conn := range clients {
		if conn != excludeConn {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("广播消息失败: %v", err)
				// 删除失败的连接
				delete(clients, conn)
			}
		}
	}
}

// broadcastLockStatus 广播锁状态给所有客户端
func broadcastLockStatus() {
	lockMsg := LockStatus{
		Type:     "lock",
		Locked:   isLocked,
		LockedBy: lockedBy,
	}
	
	msgBytes, err := json.Marshal(lockMsg)
	if err != nil {
		log.Printf("序列化锁状态失败: %v", err)
		return
	}
	
	broadcast(msgBytes, nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	
	// 为新连接分配临时客户端ID
	tempClientID := "temp_" + conn.RemoteAddr().String()
	
	mutex.Lock()
	clients[conn] = tempClientID
	mutex.Unlock()
	
	defer func() {
		mutex.Lock()
		clientID := clients[conn]
		delete(clients, conn)
		// 如果断开的客户端持有锁，自动释放锁
		if isLocked && lockedBy == clientID {
			isLocked = false
			lockedBy = ""
			log.Printf("客户端 %s 断开连接，自动释放锁", clientID)
			mutex.Unlock()
			broadcastLockStatus()
		} else {
			mutex.Unlock()
		}
		conn.Close()
		log.Printf("客户端断开 | 在线: %d\n", len(clients))
	}()

	log.Printf("新客户端连接 | 在线: %d\n", len(clients))

	// 发送当前内容给新客户端
	contentMsg := Message{
		Type:    "content",
		Content: content,
	}
	contentBytes, _ := json.Marshal(contentMsg)
	if err := conn.WriteMessage(websocket.TextMessage, contentBytes); err != nil {
		log.Printf("发送内容失败: %v", err)
		return
	}

	// 发送当前锁状态给新客户端
	lockMsg := LockStatus{
		Type:     "lock",
		Locked:   isLocked,
		LockedBy: lockedBy,
	}
	lockBytes, _ := json.Marshal(lockMsg)
	if err := conn.WriteMessage(websocket.TextMessage, lockBytes); err != nil {
		log.Printf("发送锁状态失败: %v", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			// 兼容旧格式
			mutex.Lock()
			content = string(msg)
			contentMsg := Message{
				Type:    "content",
				Content: content,
			}
			contentBytes, _ := json.Marshal(contentMsg)
			for client := range clients {
				client.WriteMessage(websocket.TextMessage, contentBytes)
			}
			mutex.Unlock()
			continue
		}

		mutex.Lock()
		clients[conn] = message.ClientID

		switch message.Type {
		case "lock":
			if !isLocked {
				isLocked = true
				lockedBy = message.ClientID
				log.Printf("客户端 %s 抢到锁", message.ClientID)
				mutex.Unlock()
				broadcastLockStatus()
			} else {
				mutex.Unlock()
			}

		case "unlock":
			if isLocked && lockedBy == message.ClientID {
				isLocked = false
				lockedBy = ""
				log.Printf("客户端 %s 释放锁", message.ClientID)
				mutex.Unlock()
				broadcastLockStatus()
			} else {
				mutex.Unlock()
			}

		case "content":
			if isLocked && lockedBy == message.ClientID {
				content = message.Content
				contentMsg := Message{
					Type:    "content",
					Content: content,
				}
				contentBytes, _ := json.Marshal(contentMsg)
				mutex.Unlock()
				// 广播给除了发送者之外的所有客户端
				broadcast(contentBytes, conn)
			} else {
				mutex.Unlock()
			}
		}
	}
}

func main() {
	port := flag.String("port", "8080", "服务端口")
	flag.Parse()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Printf("✅ 共享编辑服务启动: http://localhost:%s\n", *port)
	fmt.Printf("✅ 自动适配 PC / 手机\n")
	fmt.Printf("✅ 断开自动 30s 重连\n")
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
