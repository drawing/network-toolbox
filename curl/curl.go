package main

// https://curlconverter.com/go/
// https://github.com/quic-go/quic-go/blob/master/example/client/main.go

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/logging"
	"golang.org/x/net/http2"
)

var (
	insecure bool
	head     bool

	h3 bool
	h2 bool

	url   string
	data  string
	delay int
)

func init() {
	flag.BoolVar(&insecure, "k", false, "Ignore certificate verification")
	flag.BoolVar(&head, "I", false, "Fetch the headers only!")

	flag.BoolVar(&h3, "http3", false, "Use http3 protocol")
	flag.BoolVar(&h2, "http2", false, "Use http2 protocol")
	flag.StringVar(&url, "url", "", "Specify a URL to fetch")
	flag.StringVar(&data, "d", "", "Post data")

	flag.IntVar(&delay, "delay", 0, "Post data delay to send(ms)")
}

type reqBodyReader struct {
	body   []byte
	offset int
}

func (r *reqBodyReader) Read(p []byte) (n int, err error) {
	if r.offset == 0 && delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	if r.offset == len(r.body) {
		return 0, io.EOF
	}
	b := r.body[r.offset:]
	copy(p, b)
	if len(p) > len(b) {
		r.offset += len(b)
		// fmt.Println("read done:", len(b))
		return len(b), nil
	} else {
		r.offset += len(p)
		return len(p), nil
	}
}
func (r *reqBodyReader) Close() error {
	return nil
}

func main() {
	flag.Parse()

	if h3 && h2 {
		panic("-h3 and -h2 cannot exist at the same time")
	}

	if flag.NArg() >= 1 {
		url = flag.Arg(0)
	}

	fmt.Println("REQUEST:", url, ", Data=", data)

	var tr http.RoundTripper

	if h3 {
		var qconf quic.Config
		// if *enableQlog {
		qconf.Tracer = func(ctx context.Context, p logging.Perspective, connID quic.ConnectionID) logging.ConnectionTracer {
			filename := fmt.Sprintf("client_%x.qlog", connID)
			log.Printf("Creating qlog file %s.\n", filename)
			// TODO
			return nil
		}

		tr = &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				// RootCAs:            pool,
				InsecureSkipVerify: insecure,
				// KeyLogWriter:       keyLog,
				// TODO: set servername
				// ServerName: "default",
			},
			QuicConfig: &qconf,
		}
	} else if h2 {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}
		http2.ConfigureTransport(transport)
		tr = transport
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}
	}

	method := "GET"
	if head {
		method = "HEAD"
	}
	if len(data) > 0 {
		method = "POST"
	}

	client := &http.Client{
		Transport: tr,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	if len(data) > 0 {
		req.Body = &reqBodyReader{[]byte(data), 0}
		// req.ContentLength = int64(len(data))
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

	if h3 {
		h3Tr, ok := tr.(*http3.RoundTripper)
		if ok {
			h3Tr.Close()
		}
	}
}
