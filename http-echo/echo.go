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

type httpHandler struct {
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const defaultContent = "echo ok"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (ih *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")

	responseBody := defaultContent

	switch content {
	case "echo":
		arg_data := r.FormValue("data")
		responseBody = arg_data
	case "random":
		length, err := strconv.Atoi(r.FormValue("length"))
		if err != nil {
			responseBody = "random length parse failed " + err.Error()
			break
		}
		responseBody = randString(length)
	default:
	}

	speed, err := strconv.Atoi(r.FormValue("speed"))
	if err != nil || speed == 0 {
		io.WriteString(w, responseBody)
	} else {
		for offset := 0; offset < len(responseBody); offset += speed {
			end := offset + speed
			if end > len(responseBody) {
				end = len(responseBody)
			}
			io.WriteString(w, responseBody[offset:end])
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(time.Second)
		}
	}
}

func main() {
	var host string
	var port int

	flag.StringVar(&host, "host", "0.0.0.0", "listen host")
	flag.IntVar(&port, "port", 80, "listen port")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Println("start listen ", addr)

	http.Handle("/", &httpHandler{})
	http.ListenAndServe(addr, nil)
}
