package main

import (
	"fmt"
	"net"
	"regexp"
)

func main() {
	serve()
}

func serve() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		panic("Failed to bind to port 4221")
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			message := fmt.Errorf("error accepting connection: %v", err.Error())
			panic(message)
		}

		go func() {
			defer conn.Close()

			req, err := ParseRequest(conn)
			if err != nil {
				fmt.Printf("error parsing request: %v", err.Error())
				return
			}

			hc := newHttpConnection(conn, req)
			handleRouting(hc)
		}()
	}
}

func handleRouting(hc HttpConnection) {
	// TODO: Refactor
	method := hc.req.Method
	path := hc.req.Path

	if method == "GET" {

		if path == "/" {
			hc.Respond()
			return
		}

		pattern := regexp.MustCompile(`/echo/([^/]+)`)
		matches := pattern.FindStringSubmatch(path)
		if len(matches) == 2 {
			hc.res.Body = matches[1]
			hc.Respond()
			return
		}

	}

	hc.res.Status = STATUS_NOT_FOUND
	hc.Respond()
}
