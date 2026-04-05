package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// EchoHandler 处理 HTTP 请求并根据参数返回不同的响应
type EchoHandler struct{}

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultContent = "echo ok"
)

// init 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

// generateRandomString 生成指定长度的随机字符串
func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// getResponseContent 根据请求参数确定响应内容
func getResponseContent(contentType, data, lengthStr string) string {
	switch contentType {
	case "echo":
		return data
	case "random":
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return fmt.Sprintf("random length parse failed: %v", err)
		}
		return generateRandomString(length)
	default:
		return defaultContent
	}
}

// writeResponseWithSpeed 按照指定速度分块写入响应
func writeResponseWithSpeed(w http.ResponseWriter, responseBody string, speed int) {
	for offset := 0; offset < len(responseBody); offset += speed {
		end := offset + speed
		if end > len(responseBody) {
			end = len(responseBody)
		}
		
		if _, err := io.WriteString(w, responseBody[offset:end]); err != nil {
			log.Printf("Error writing response: %v", err)
			return
		}
		
		// 尝试刷新响应缓冲区
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		
		time.Sleep(time.Second)
	}
}

// ServeHTTP 处理 HTTP 请求
func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.FormValue("content")
	data := r.FormValue("data")
	lengthStr := r.FormValue("length")
	speedStr := r.FormValue("speed")

	// 确定响应内容
	responseBody := getResponseContent(contentType, data, lengthStr)

	// 检查是否需要按速度限制发送
	if speedStr != "" {
		speed, err := strconv.Atoi(speedStr)
		if err == nil && speed > 0 {
			writeResponseWithSpeed(w, responseBody, speed)
			return
		}
	}

	// 直接发送响应
	if _, err := io.WriteString(w, responseBody); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	var host string
	var port int

	flag.StringVar(&host, "host", "0.0.0.0", "监听主机地址")
	flag.IntVar(&port, "port", 80, "监听端口")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("启动 HTTP Echo 服务，监听地址: %s", addr)

	http.Handle("/", &EchoHandler{})
	log.Fatal(http.ListenAndServe(addr, nil))
}
