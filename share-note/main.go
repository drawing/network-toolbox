package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	content string
	clients = make(map[*websocket.Conn]bool)
	mutex   sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const html = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>共享实时输入框</title>
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial, sans-serif;
            background: #f7f8fa;
            padding: 10px;
        }

        .container {
            max-width: 900px;
            margin: 0 auto;
        }

        .status {
            padding: 8px 12px;
            border-radius: 8px;
            font-size: 14px;
            margin-bottom: 10px;
            display: none;
        }

        .status-connected {
            background: #e6f7e6;
            color: #2a962a;
            display: block;
        }

        .status-disconnected {
            background: #fff1f0;
            color: #d32f2f;
            display: block;
        }

        textarea {
            width: 100%;
            padding: 15px;
            font-size: 16px;
            border: 1px solid #eee;
            border-radius: 12px;
            background: #fff;
            resize: none;
            outline: none;
            box-shadow: 0 2px 8px rgba(0,0,0,0.05);
        }

        /* PC 样式 */
        @media screen and (min-width: 768px) {
            textarea {
                height: 500px;
            }
        }

        /* 手机样式（自动适配）*/
        @media screen and (max-width: 767px) {
            body {
                padding: 5px;
                background: #fff;
            }
            textarea {
                height: 82vh;
                font-size: 18px;
                border-radius: 8px;
                border: none;
                box-shadow: none;
            }
            .status {
                font-size: 12px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div id="status" class="status">连接中...</div>
        <textarea id="editor" placeholder="输入内容，实时多人同步…"></textarea>
    </div>

    <script>
        const editor = document.getElementById('editor');
        const status = document.getElementById('status');
        let ws;
        let reconnectInterval = 30000; // 30秒重连
        let isConnected = false;

        // 自动连接
        function connect() {
            // 关闭旧连接
            if (ws) {
                ws.close();
            }

            // 新建 WebSocket
            ws = new WebSocket("ws://" + location.host + "/ws");

            // 连接成功
            ws.onopen = function () {
                isConnected = true;
                status.className = "status status-connected";
                status.innerText = "✅ 已连接 · 实时同步中";
            };

            // 收到消息
            ws.onmessage = function (evt) {
                editor.value = evt.data;
            };

            // 连接关闭/失败
            ws.onclose = function () {
                if (isConnected) {
                    isConnected = false;
                    status.className = "status status-disconnected";
                    status.innerText = "🔌 已断开，" + (reconnectInterval / 1000) + "秒后自动重试...";
                }
                // 定时重连
                setTimeout(connect, reconnectInterval);
            };

            // 连接错误
            ws.onerror = function (err) {
                console.log("WebSocket 错误", err);
                ws.close();
            };
        }

        // 输入变化 → 发送
        editor.oninput = function () {
            if (ws && isConnected) {
                try {
                    ws.send(editor.value);
                } catch (e) {}
            }
        };

        // 启动
        connect();
    </script>
</body>
</html>
`

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.New("index").Parse(html)
	tpl.Execute(w, nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer func() {
		conn.Close()
	}()

	mutex.Lock()
	clients[conn] = true
	conn.WriteMessage(websocket.TextMessage, []byte(content))
	mutex.Unlock()

	log.Printf("新客户端连接 | 在线: %d\n", len(clients))

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		mutex.Lock()
		content = string(msg)
		for client := range clients {
			client.WriteMessage(websocket.TextMessage, msg)
		}
		mutex.Unlock()
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	log.Printf("客户端断开 | 在线: %d\n", len(clients))
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
