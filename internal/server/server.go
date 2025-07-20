package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

func ListenAndServe(addr string, handler Handler) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnction(conn, handler)
	}
}

func handleConnction(conn net.Conn, handler Handler) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		req, err := ReadRequest(reader)
		if err != nil {
			fmt.Fprintf(conn, "HTTP/1.1 400 Bad Request\r\n\r\n%s\n", err)
			return
		}
		if req == nil {
			return
		}

		w := &simpleResponseWriter{conn: conn}
		handler.Serve(w, req)
		w.Flush()

		close := strings.ToLower(req.Headers["Connection"]) == "close"
		if close || req.Protocol == "HTTP/1.0" {
			return
		}
	}
}

type route struct {
	pattern *regexp.Regexp
	handler Handler
}
type ServeMux struct {
	routes []route
}

func NewServeMux() *ServeMux {
	return &ServeMux{}
}

func (m *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	converted := convertPattern(pattern)
	re := regexp.MustCompile(converted)
	m.routes = append(m.routes, route{pattern: re, handler: HandlerFunc(handler)})
}

func (m *ServeMux) Serve(w ResponseWriter, r *Request) {
	for _, rt := range m.routes {
		if rt.pattern.MatchString(r.Path) {
			matches := rt.pattern.FindStringSubmatch(r.Path)
			params := make(map[string]string)
			for i, name := range rt.pattern.SubexpNames() {
				if i > 0 && name != "" {
					params[name] = matches[i]
				}
			}
			r.Params = params
			rt.handler.Serve(w, r)
			return
		}
	}
	fmt.Fprintln(w, "HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\nNot Found")
}

func convertPattern(p string) string {
	// Replace ":param" with regex group
	re := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
	return "^" + re.ReplaceAllString(p, `(?P<$1>[^/]+)`) + "$"
}

func parseQuery(qs string) map[string]string {
	params := make(map[string]string)
	if qs == "" {
		return params
	}
	pairs := strings.SplitSeq(qs, "&")
	for pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		}
	}
	return params
}
