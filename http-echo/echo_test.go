package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestEchoHandler_EchoContent 测试 echo 内容处理
func TestEchoHandler_EchoContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?content=echo&data=abc", nil)
	w := httptest.NewRecorder()

	(&EchoHandler{}).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if string(data) != "abc" {
		t.Errorf("expected 'abc', got '%v'", string(data))
	}
}

// TestEchoHandler_DefaultContent 测试默认内容处理
func TestEchoHandler_DefaultContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?content=empty", nil)
	w := httptest.NewRecorder()

	(&EchoHandler{}).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if string(data) != defaultContent {
		t.Errorf("expected '%v', got '%v'", defaultContent, string(data))
	}
}

// TestEchoHandler_RandomContent 测试随机内容处理
func TestEchoHandler_RandomContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?content=random&length=1234", nil)
	w := httptest.NewRecorder()

	(&EchoHandler{}).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if len(data) != 1234 {
		t.Errorf("expected length %d, got %d", 1234, len(data))
	}
}

// TestEchoHandler_SpeedContent 测试带速度限制的内容处理
func TestEchoHandler_SpeedContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?content=random&length=10&speed=4", nil)
	w := httptest.NewRecorder()

	handler := &EchoHandler{}
	go handler.ServeHTTP(w, req)

	// 等待一小段时间让服务器开始处理
	time.Sleep(500 * time.Millisecond)

	expectedChunkLengths := []int{4, 4, 2}

	for i, expectedLen := range expectedChunkLengths {
		data := make([]byte, 20)
		n, err := w.Body.Read(data)
		if err != nil {
			t.Fatalf("error reading response body at chunk %d: %v", i, err)
		}
		if n != expectedLen {
			t.Errorf("chunk %d: expected length %d, got %d", i, expectedLen, n)
		}
		// 等待下一个 chunk
		time.Sleep(time.Second)
	}
}
