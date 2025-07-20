# MiniHTTP Server

A lightweight HTTP server built from scratch in Go using the net package. It handles basic routing, dynamic paths, query parameters, and concurrent request handling using goroutines.

## Features

- **Routing**: Matches HTTP requests to handlers based on URL patterns.
- **Dynamic Paths**: Captures path segments (e.g., /user/:name extracts name).
- **Query Parameters**: Parses query strings (e.g., /search?q=foo).
- **POST Body Handling**: Processes request bodies up to 64 KiB.
- **Custom Headers**: Supports adding headers like X-Foo: bar.
- **HTTP/1.1 Keep-Alive**: Handles persistent connections for multiple requests.
- **Concurrent Connections**: Uses goroutines for efficient request handling.

## How It Uses the `net` Package

MiniHTTP leverages Go’s `net` package to manage TCP connections at a low level:

1. **TCP Listener**:
   - `net.Listen("tcp", ":8080")` creates a TCP listener on port 8080 to accept incoming connections.

2. **Connection Handling**:
   - `l.Accept()` grabs incoming TCP connections (`net.Conn`), spawning a goroutine per connection for concurrency.
   - Each connection is processed in `handleConnction`, which reads requests and writes responses.

3. **Reading Requests**:
   - A `bufio.Reader` on `net.Conn` reads raw HTTP request data (method, path, headers, body).
   - `ReadRequest` parses the TCP stream into a `Request` struct, handling query params and body (up to 64 KiB).

4. **Writing Responses**:
   - The `simpleResponseWriter` uses `net.Conn.Write` to send HTTP responses (status, headers, body).
   - Responses are formatted  (e.g., `HTTP/1.1 200 OK\r\n`) and written to the TCP connection.

5. **Keep-Alive Support**:
   - Supports HTTP/1.1 persistent connections by checking the `Connection` header, looping to handle multiple requests unless closed or HTTP/1.0 is used.

## Usage

1. **Run**:

   ```bash
   go run main.go
   ```

   Server starts on `:8080`.

2. **Test Endpoints**:
   - `curl http://localhost:8080/` → `Hello server`
   - `curl http://localhost:8080/user/Alice` → `Hello, Alice`
   - `curl http://localhost:8080/search?q=foo` → `Search for: foo`
   - `curl -X POST -d "test" http://localhost:8080/submit` → `Received: test`

