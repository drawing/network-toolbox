# curl

A simple HTTP client tool supporting HTTP/1.1, HTTP/2 and HTTP/3 protocols, similar to the curl command-line tool.

## Features

- 📡 **Multi-protocol support**: Supports HTTP/1.1, HTTP/2 and HTTP/3 protocols
- 🔒 **Certificate verification**: Optional SSL certificate verification
- 📤 **POST requests**: Supports sending POST data
- ⏱️ **Delayed sending**: Supports setting delay before sending POST data
- 📋 **Header requests**: Supports fetching only HTTP headers
- 🖨️ **Detailed output**: Shows response status code and body

## Technical Implementation

- **Language**: Go
- **HTTP client**: Standard library `net/http`
- **HTTP/2 support**: `golang.org/x/net/http2`
- **HTTP/3 support**: `github.com/quic-go/quic-go`

## Quick Start

### Prerequisites

- Go 1.16 or higher
- Git

### Installation and Running

1. **Clone the project**

```bash
git clone https://github.com/yourusername/network-toolbox.git
cd network-toolbox/curl
```

2. **Install dependencies**

```bash
go get github.com/quic-go/quic-go
go get golang.org/x/net/http2
```

3. **Build and run**

```bash
# Build
go build -o curl .

# Run
./curl https://example.com
```

## Usage

### Basic Usage

```bash
# Fetch webpage using default HTTP/1.1 protocol
./curl https://example.com

# Specify URL using -url parameter
./curl -url https://example.com
```

### Protocol Selection

```bash
# Use HTTP/2 protocol
./curl -http2 https://example.com

# Use HTTP/3 protocol
./curl -http3 https://example.com
```

### POST Requests

```bash
# Send POST data
./curl -d "key=value&another=value" https://example.com/api

# Send POST data with delay (1000ms)
./curl -d "key=value" -delay 1000 https://example.com/api
```

### Other Options

```bash
# Ignore SSL certificate verification
./curl -k https://self-signed.example.com

# Fetch only headers
./curl -I https://example.com
```

## Command-line Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `-k` | Ignore SSL certificate verification | `false` |
| `-I` | Fetch only headers | `false` |
| `-http3` | Use HTTP/3 protocol | `false` |
| `-http2` | Use HTTP/2 protocol | `false` |
| `-url` | Specify URL to fetch | `""` |
| `-d` | POST data | `""` |
| `-delay` | Delay before sending POST data (milliseconds) | `0` |

## Examples

### 1. Basic GET Request

```bash
./curl https://example.com
```

Output:
```
REQUEST: https://example.com, Data=
Response status: 200 OK
Response body:
<!doctype html>
<html>
<head>
    <title>Example Domain</title>
    ...
</html>
```

### 2. POST Request

```bash
./curl -d "name=test&value=123" https://httpbin.org/post
```

Output:
```
REQUEST: https://httpbin.org/post, Data=name=test&value=123
Response status: 200 OK
Response body:
{
  "args": {},
  "data": "",
  "files": {},
  "form": {
    "name": "test",
    "value": "123"
  },
  ...
}
```

### 3. Using HTTP/3 Protocol

```bash
./curl -http3 https://www.google.com
```

Output:
```
REQUEST: https://www.google.com, Data=
Response status: 200 OK
Response body:
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Google</title>
    ...
</html>
```

## Project Structure

```
curl/
├── curl.go          # Main program file
└── README.md        # Project documentation
```

## Core Functionality Implementation

### Protocol Selection

- Selects HTTP/1.1, HTTP/2 or HTTP/3 protocol based on command-line parameters
- Uses quic-go library for HTTP/3 QUIC protocol implementation

### Delayed Sending

- Implements `delayedBodyReader` struct to support delay before sending POST data
- Delay is only applied on first read

### Error Handling

- Provides detailed error messages for URL validation, transport creation and request sending
- Adds HTTP client timeout to prevent hanging

## Notes

- HTTP/3 functionality depends on quic-go library and may require additional system dependencies
- Currently, HTTP/3 qlog functionality only creates log file path, full logging not implemented
- For HTTPS requests, SSL certificate verification is enabled by default, use `-k` parameter to ignore

## Extension Suggestions

- Add support for more HTTP methods (PUT, DELETE, etc.)
- Implement file upload functionality
- Add HTTP header customization
- Implement proxy support
- Add more output format options (JSON, quiet mode, etc.)

## License

MIT License

## Contributing

Welcome to submit Issues and Pull Requests to improve this project!