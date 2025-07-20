package server

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

const maxBodyBytes = 64 * 1024 // 64 KiB

var (
	ErrBadRequestLine   = errors.New("malformed request line")
	ErrHeaderTooLarge   = errors.New("header section too large")
	ErrBadContentLength = errors.New("invalid Content-Length")
	ErrBodyTooLarge     = errors.New("request body too large")
	ErrUnexpectedEOF    = errors.New("unexpected EOF")
)

type Request struct {
	Headers  map[string]string
	Params   map[string]string
	Query    map[string]string
	Method   string
	Path     string
	Protocol string
	Body     string
}

func ReadRequest(reader *bufio.Reader) (*Request, error) {
	// Request line
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, ErrBadRequestLine
	}

	method, rawURL, protocol := parts[0], parts[1], parts[2]

	path := rawURL
	queryString := ""
	if idx := strings.Index(rawURL, "?"); idx != -1 {
		path = rawURL[:idx]
		queryString = rawURL[idx+1:]
	}

	// Headers
	headers := make(map[string]string)
	for {
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		if line == "" { // empty line => end of headers
			break
		}
		colon := strings.IndexByte(line, ':')
		if colon == -1 {
			continue // ignore malformed header lines
		}
		key := strings.TrimSpace(line[:colon])
		val := strings.TrimSpace(line[colon+1:])
		headers[key] = val
	}

	// Body
	var body string
	if clStr, ok := headers["Content-Length"]; ok {
		cl, err := strconv.Atoi(clStr)
		if err != nil || cl < 0 {
			return nil, ErrBadContentLength
		}
		if cl > maxBodyBytes {
			return nil, ErrBodyTooLarge
		}
		buf := make([]byte, cl)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, ErrUnexpectedEOF
		}
		body = string(buf)
	}

	return &Request{
		Method:   method,
		Path:     path,
		Protocol: protocol,
		Headers:  headers,
		Body:     body,
		Query:    parseQuery(queryString),
	}, nil
}

func readLine(r *bufio.Reader) (string, error) {
	const maxLine = 4096 // per RFC 7230 recommendation
	s, err := r.ReadString('\n')
	if err != nil {
		if err == io.EOF && len(s) == 0 {
			return "", err
		}
		return "", ErrUnexpectedEOF
	}
	if len(s) > maxLine {
		return "", ErrHeaderTooLarge
	}
	return strings.TrimRight(s, "\r\n"), nil
}
