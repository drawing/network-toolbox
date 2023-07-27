package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEchoContentHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?content=echo&data=abc", nil)
	w := httptest.NewRecorder()

	(&httpHandler{}).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "abc" {
		t.Errorf("expected abc got %v", string(data))
	}
}

func TestDefaultContentHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?content=empty", nil)
	w := httptest.NewRecorder()

	(&httpHandler{}).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != defaultContent {
		t.Errorf("expected %v got %v", defaultContent, string(data))
	}
}

func TestRandomContentHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?content=random&length=1234", nil)
	w := httptest.NewRecorder()

	(&httpHandler{}).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if len(data) != 1234 {
		t.Errorf("expected len %d got %d", 1234, len(data))
	}
}

func TestSpeedContentHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?content=random&length=10&speed=4", nil)
	w := httptest.NewRecorder()

	handler := &httpHandler{}
	go handler.ServeHTTP(w, req)

	time.Sleep(time.Microsecond * 500)

	expectLens := []int{4, 4, 2}

	for i := 0; i < len(expectLens); i++ {
		data := make([]byte, 20)
		n, err := w.Body.Read(data)
		if err != nil {
			break
		}
		if n != expectLens[i] {
			t.Errorf("expected len %d got %d", expectLens[i], n)
		}
		time.Sleep(time.Second)
	}
}
