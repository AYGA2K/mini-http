package server

import (
	"fmt"
	"net"
	"strconv"
)

type ResponseWriter interface {
	Write([]byte) (int, error)
	Flush() error
	Header(string, string)
	Body([]byte)
	Send()
	StatusCode(int)
}

type simpleResponseWriter struct {
	conn       net.Conn
	headers    map[string]string
	body       []byte
	statusCode int
	headerSent bool
}

func (w *simpleResponseWriter) StatusCode(code int) {
	w.statusCode = code
}

func (w *simpleResponseWriter) Header(key, value string) {
	if w.headers == nil {
		w.headers = make(map[string]string)
	}
	w.headers[key] = value
}

func (w *simpleResponseWriter) Body(body []byte) {
	w.body = body
}

func (w *simpleResponseWriter) Write(data []byte) (int, error) {
	return w.conn.Write(data)
}

func (w *simpleResponseWriter) Flush() error {
	return nil
}

func (w *simpleResponseWriter) Send() {
	if w.statusCode == 0 {
		w.statusCode = 200
	}
	if w.headers == nil {
		w.headers = make(map[string]string)
	}
	// defaults
	if w.headers["Content-Type"] == "" {
		w.headers["Content-Type"] = "text/plain"
	}
	w.headers["Content-Length"] = strconv.Itoa(len(w.body))

	fmt.Fprintf(w.conn, "HTTP/1.1 %d %s\r\n", w.statusCode, statusText(w.statusCode))
	for k, v := range w.headers {
		fmt.Fprintf(w.conn, "%s: %s\r\n", k, v)
	}
	fmt.Fprint(w.conn, "\r\n")

	fmt.Fprint(w.conn, string(w.body))
}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	default:
		return "Unknown"
	}
}
