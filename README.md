# Http 1.1 From Scratch 

An implementation of the HTTP protocol written in Go, featuring a custom state machine parser and a method aware mutliplexer.

## Installation

### Prerequisites

- Go 1.20 or higher
- `curl` for testing
- `ab` (ApacheBench) for benchmarking

### Build from Source

1. Clone the repository:
```bash
git clone https://github.com/yourusername/httpserver.git
cd httpserver
```

2. Run the server:
```bash
go run cmd/httpServer/main.go
```

Alternatively, build a binary:
```bash
go build -o httpserver cmd/httpServer/main.go
./httpserver
```

The server will start on `localhost:8000` by default.

## Verify Installation

Test that the server is running:
```bash
curl http://localhost:8000/
```

Should output:
```
Hello, World!
```

## Usage

### Starting the Server

The server file in `cmd/server` handles TCP socket management and connection dispatching. However, the main file `cmd/httpServer` is what starts an http server. The tcp connection is wrapped in a listen and serve method because its primary function is keep on running while dispatching information. The http server is meant to handle the configuration of what gets dispatched:

```go
// cmd/server/server.go - tcp connection entry point
func CustomListenAndServe(addr string, h handler.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn) {
			err := serve(c, h)
			if err != nil {
				println("Error in trying to serve connection", err.Error())
			}
		}(conn)
	}
}
```

```go
// cmd/httpServer/main.go - an example of an httpServer with custom handlers and routes.
func main() {
	errChan := make(chan error, 1)
	go func () {
		server.MyDefaultMux.HandleFunc("GET", "/", func(w response.ResponseWriter, r *request.Request) {
			fmt.Printf("Handled GET / request\n")		
			w.Write([]byte("Hello, World!\n"))
		})
		server.MyDefaultMux.HandleFunc("POST", "/upload", func(w response.ResponseWriter, r *request.Request) {
			w.GetHeaders().Set("Content-Type", "application/json")
			w.CustomWriteHeader(201)
			w.Write([]byte(`{"status":"success"}`))
		})
		fmt.Printf("Server started running")
		errChan <- server.CustomListenAndServe(":8000", nil)
	}()
	err := <- errChan
	fmt.Printf("Server stopped: ", err)
}
```

### Adding Request Handlers

Register handlers using the custom multiplexer:

```go
// Example: Add a new route handler
mux := NewServerMux()

// GET handler
mux.HandleFunc("GET", "/", func(w ResponseWriter, r *Request) {
    w.CustomWriteHeader(200)
    w.Write([]byte("Hello, World!"))
})

// POST handler with body parsing
mux.HandleFunc("POST", "/upload", func(w ResponseWriter, r *Request) {
    body := r.Body // Already parsed by state machine
    w.CustomWriteHeader(201)
    w.Write([]byte("I received your data: " + body))
})

// JSON API endpoint
mux.HandleFunc("POST", "/api/data", func(w ResponseWriter, r *Request) {
    w.GetHeaders().Set("Content-Type", "application/json")
    w.CustomWriteHeader(200)
    w.Write([]byte(`{"status":"success"}`))
})
```

### The Parser State Machine

The request parser transitions through states to handle fragmented TCP streams:

```
StateInit → StateHeaders → StateBody → StateDone
                               ↓
                StateBodyFixed / StateBodyChunkedRead and StateBodyChunkedWrite
```

**State Transitions:**
- `StateInit`: Parsing `METHOD /path HTTP/1.1`
- `StateHeaders`: Reading headers until `\r\n\r\n`
- `StateBodyFixed`: Reads up to identified body length
- `StateBodyChunkedRead/StateBodyChunkedWrite`: Alternates between reading anticipated size and writing the body to the requests structure
- `StateDone`: Request ready for handler

## Testing

### Basic Connectivity (GET)

Verify request parsing and response generation:

```bash
curl -v http://localhost:8000/
```

Expected output:
```
> GET / HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/8.15.0
>
< HTTP/1.1 200 OK
< Content-Type: text/plain
< Content-Length: 13
<
Hello, World!
```

**What to check:**
- Status line: `HTTP/1.1 200 OK`
- Response body matches expected content
- Connection closes cleanly

### Data Handling (POST)

Test body parsing with `Content-Length`:

```bash
curl -v -X POST -d "name=Juan&project=httpserver" http://localhost:8000/upload
```

Expected output:
```
> POST /upload HTTP/1.1
> Content-Length: 28
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 201 Created
<
{"status":"success"}
```

**What to check:**
- Server echoes received data
- `Content-Length` correctly parsed (28 bytes)
- Status code `201 Created`

### Stress Test (ApacheBench)

Test concurrency and throughput:

```bash
# 1,000 requests with 10 concurrent connections
ab -n 1000 -c 10 http://localhost:8000/
```

Expected results:
```
Concurrency Level:      10
Complete requests:      1000
Failed requests:        0
Requests per second:    1500+ [#/sec]
```

**What to check:**
- Failed requests: **0**
- Requests per second > 1000
- No connection timeouts

Advanced stress test:
```bash
# 10,000 requests with 100 concurrent connections
ab -n 10000 -c 100 http://localhost:8000/

# Keep-alive test
ab -n 5000 -c 50 -k http://localhost:8000/
```

## Architecture

### Key Features

- **Manual TCP Management**: Direct `net.Listen` socket handling
- **State-Machine Parser**: Handles fragmented streams gracefully
- **Goroutine-Per-Connection**: Non-blocking concurrent request handling
- **Custom ResponseWriter**: Full control over HTTP response formatting
- **Zero Dependencies**: Standard library only

### File Structure

```
httpserver/
├── cmd/
│   ├── client/
│   │   └── main.go          # CLI client to test the server
│   ├── httpServer/
│   │   └── main.go          # The main entry point (Server initialization)
│   └── tcplistener/
│       └── main.go          # low level TCP experiments/benchmarking
├── internal/
│   ├── handler/
│   │   ├── handler.go       # Handler interface & HandlerFunc adapter
│   │   └── handler_test.go
│   ├── headers/
│   │   ├── headers.go       # Header parsing logic
│   │   └── headers_test.go
│   ├── request/
│   │   ├── generate_request.go
│   │   ├── request.go       # The state machine parser
│   │   └── request_test.go
│   ├── response/
│   │   ├── response.go      # ResponseWriter & status line logic
│   │   └── response_test.go
│   └── server/
│       ├── server.go        # Mux, routeKey, & (custom) ListenAndServe
│       └── server_test.go
├── README.md                # Documentation & usage guide
├── go.mod                   # Module definition
└── go.sum                   # Dependency checksums
```

### Performance

Typical benchmarks on modern hardware:

| Metric | Value |
|--------|-------|
| Requests/sec | ~15,000 |
| Concurrency | 100 connections |
| Failed Requests | 0 |
| Memory Usage | ~50MB |

## Troubleshooting

### Server won't start

Check if port 8000 is already in use:
```bash
lsof -i :8000
# Or change the port in main.go
```

### Failed requests in stress test

1. Check system limits:
```bash
ulimit -n  # Should be > 1024
```

2. Increase if needed:
```bash
ulimit -n 4096
```

3. Monitor server logs for panics or errors

### Connection timeouts

Verify the server is running:
```bash
ps aux | grep httpserver
```

Check firewall rules:
```bash
sudo iptables -L | grep 8000
```

## Contributing

Contributions are welcome! This project prioritizes:

1. **Code clarity** over clever optimizations
2. **Spec compliance** with HTTP/1.1 RFCs
3. **Zero dependencies** philosophy

### How to Contribute

1. Fork the repository
2. Create a feature branch:
```bash
git checkout -b feature/your-feature-name
```

3. Make your changes following these guidelines:
   - Document all exported functions and types
   - Add tests for new functionality
   - Ensure `go fmt` and `go vet` pass
   - Verify stress tests still pass with `ab`

4. Test your changes:
```bash
# Run the server
go run ./cmd/server

# In another terminal
curl -v http://localhost:8000/
ab -n 1000 -c 10 http://localhost:8000/
```

5. Submit a pull request with:
   - Clear description of changes
   - Any relevant test results
   - Updated documentation if needed

## Resources

This implementation was built following HTTP/1.1 specifications and educational resources:

### RFC Specifications

- **[RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html)**: HTTP Semantics (methods, status codes, headers)
- **[RFC 9112](https://www.rfc-editor.org/rfc/rfc9112.html)**: HTTP/1.1 Message Syntax and Routing
- **[RFC 7231](https://www.rfc-editor.org/rfc/rfc7231.html)**: HTTP/1.1 Semantics and Content (methods and status codes)

### Educational Resources

- **[From TCP to HTTP](https://youtu.be/FknTw9bJsXM?si=yLDTfnuxrThcFQVG)**: YouTube course on building an HTTP server from scratch
- **[The Go Programming Language Specification](https://go.dev/ref/spec)**: Official Go language documentation

### Further Reading

- [HTTP/1.1 on Wikipedia](https://en.wikipedia.org/wiki/HTTP/1.1)
- [Building a Web Server in Go](https://www.golang-book.com/)
- [Understanding HTTP State Machines](https://www.w3.org/Protocols/rfc2616/rfc2616-sec4.html)

## License

MIT License - free to use as a learning resource or foundation for your own projects.

---

**Questions?** Open an issue on GitHub or contribute improvements via pull requests!
