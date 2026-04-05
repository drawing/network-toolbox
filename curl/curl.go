package main

// curl.go - A simple HTTP client supporting HTTP/1.1, HTTP/2 and HTTP/3
// Inspired by: https://curlconverter.com/go/
// HTTP/3 example: https://github.com/quic-go/quic-go/blob/master/example/client/main.go

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
	// Command line flags
	insecure bool // Ignore certificate verification
	headOnly bool // Fetch only headers
	http3Enabled bool // Use HTTP/3 protocol
	http2Enabled bool // Use HTTP/2 protocol
	targetURL string // URL to fetch
	postData string // Data to send with POST request
	delayMs int // Delay before sending POST data (milliseconds)
)

func init() {
	flag.BoolVar(&insecure, "k", false, "Ignore certificate verification")
	flag.BoolVar(&headOnly, "I", false, "Fetch the headers only")
	flag.BoolVar(&http3Enabled, "http3", false, "Use HTTP/3 protocol")
	flag.BoolVar(&http2Enabled, "http2", false, "Use HTTP/2 protocol")
	flag.StringVar(&targetURL, "url", "", "Specify a URL to fetch")
	flag.StringVar(&postData, "d", "", "Post data")
	flag.IntVar(&delayMs, "delay", 0, "Delay before sending POST data (ms)")
}

// delayedBodyReader implements io.Reader with support for delay before first read
type delayedBodyReader struct {
	body   []byte
	offset int
	delay  time.Duration
}

// Read implements io.Reader interface with optional delay on first read
func (r *delayedBodyReader) Read(p []byte) (n int, err error) {
	// Apply delay only on first read
	if r.offset == 0 && r.delay > 0 {
		time.Sleep(r.delay)
	}
	
	// Check if we've reached the end of the body
	if r.offset >= len(r.body) {
		return 0, io.EOF
	}
	
	// Calculate how much data to read
	remaining := len(r.body) - r.offset
	readSize := len(p)
	if readSize > remaining {
		readSize = remaining
	}
	
	// Copy data to the provided buffer
	copy(p, r.body[r.offset:r.offset+readSize])
	r.offset += readSize
	
	return readSize, nil
}

// Close implements io.Closer interface
func (r *delayedBodyReader) Close() error {
	return nil
}

// createHTTPTransport creates the appropriate HTTP transport based on the protocol flag
func createHTTPTransport() (http.RoundTripper, error) {
	if http3Enabled && http2Enabled {
		return nil, fmt.Errorf("HTTP/3 and HTTP/2 cannot be enabled at the same time")
	}

	if http3Enabled {
		// Create HTTP/3 transport
		quicConfig := &quic.Config{
			Tracer: func(ctx context.Context, p logging.Perspective, connID quic.ConnectionID) logging.ConnectionTracer {
				filename := fmt.Sprintf("client_%x.qlog", connID)
				log.Printf("Creating qlog file %s.\n", filename)
				// TODO: Implement qlog file creation
				return nil
			},
		}

		return &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
			QuicConfig: quicConfig,
		}, nil
	} else if http2Enabled {
		// Create HTTP/2 transport
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}
		http2.ConfigureTransport(transport)
		return transport, nil
	} else {
		// Create default HTTP/1.1 transport
		return &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}, nil
	}
}

// determineHTTPMethod determines the HTTP method based on flags and data
func determineHTTPMethod() string {
	if headOnly {
		return "HEAD"
	}
	if len(postData) > 0 {
		return "POST"
	}
	return "GET"
}

func main() {
	flag.Parse()

	// Use first positional argument as URL if provided
	if flag.NArg() >= 1 {
		targetURL = flag.Arg(0)
	}

	// Validate URL
	if targetURL == "" {
		log.Fatal("Error: URL is required. Use -url flag or provide as positional argument.")
	}

	fmt.Printf("REQUEST: %s, Data=%s\n", targetURL, postData)

	// Create HTTP transport
	transport, err := createHTTPTransport()
	if err != nil {
		log.Fatalf("Error creating transport: %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Add timeout to prevent hanging
	}

	// Determine HTTP method
	method := determineHTTPMethod()

	// Create request
	var req *http.Request
	if len(postData) > 0 {
		// Create request with delayed body reader
		bodyReader := &delayedBodyReader{
			body:   []byte(postData),
			delay: time.Duration(delayMs) * time.Millisecond,
		}
		req, err = http.NewRequest(method, targetURL, bodyReader)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}
		// Set Content-Type header for POST requests
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest(method, targetURL, nil)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Print response
	fmt.Printf("Response status: %s\n", resp.Status)
	if !headOnly {
		fmt.Printf("Response body:\n%s\n", bodyText)
	}

	// Clean up HTTP/3 transport
	if http3Enabled {
		if h3Tr, ok := transport.(*http3.RoundTripper); ok {
			h3Tr.Close()
		}
	}
}
