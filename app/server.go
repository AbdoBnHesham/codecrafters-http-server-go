package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

const (
	HTTP_VERSION                 = "HTTP/1.1"
	STATUS_OK                    = "200 OK"
	STATUS_NOT_FOUND             = "404 Not Found"
	STATUS_INTERNAL_SERVER_ERROR = "500 Internal Server Error"
)

type Request struct {
	Headers map[string]string
	Method  string
	Path    string
}

func parseRequest(reader io.Reader) (Request, error) {
	buffer := make([]byte, 1024)
	_, err := reader.Read(buffer)
	if err != nil {
		return Request{}, fmt.Errorf("error reading response: %v", err.Error())
	}

	lines := bytes.Split(buffer, []byte("\r\n"))

	// parse status line
	sections := bytes.Split(lines[0], []byte(" "))
	if len(sections) < 2 {
		return Request{}, fmt.Errorf("error parsing status line: %s", sections)
	}
	method := sections[0]
	path := sections[1]

	// TODO parse headers
	// TODO parse body

	return Request{
		Method: string(method),
		Path:   string(path),
	}, nil
}

type HttpConnection struct {
	tcpConn net.Conn
	req     Request
}

func newHttpConnection(conn net.Conn, req Request) HttpConnection {
	return HttpConnection{tcpConn: conn, req: req}
}

func (hc *HttpConnection) Respond(status string) error {
	response := []byte(fmt.Sprintf("%s %s \r\n\r\n", HTTP_VERSION, status))
	_, err := hc.tcpConn.Write(response)
	if err != nil {
		return fmt.Errorf("error writing response: %v", err.Error())
	}
	return nil
}

type RoutesMap map[string]func(HttpConnection) error

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		panic("Failed to bind to port 4221")
	}

	routesMap := RoutesMap{
		"/": func(hc HttpConnection) error {
			hc.Respond(STATUS_OK)
			return nil
		},
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			message := fmt.Errorf("error accepting connection: %v", err.Error())
			panic(message)
		}

		go func() {
			defer conn.Close()

			req, err := parseRequest(conn)
			if err != nil {
				fmt.Printf("error parsing request: %v", err.Error())
				return
			}

			hc := newHttpConnection(conn, req)
			handleRouting(hc, routesMap)
		}()
	}
}

func handleRouting(hc HttpConnection, routesMap RoutesMap) error {
	if route := routesMap[hc.req.Path]; route != nil {
		if err := route(hc); err != nil {
			hc.Respond(STATUS_INTERNAL_SERVER_ERROR)
		}
	} else {
		hc.Respond(STATUS_NOT_FOUND)
	}

	return nil
}
