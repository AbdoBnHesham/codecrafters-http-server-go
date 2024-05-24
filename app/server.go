package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	serve()
}

func serve() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		panic("Failed to bind to port 4221")
	}

	config := newConfig()

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
			hc.config = config
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

		if path == "/user-agent" {
			for k, v := range hc.req.Headers {
				if strings.ToLower(k) == "user-agent" {
					hc.res.Body = v
					// TODO: automate this
					hc.res.Headers["Content-Type"] = "text/plain"
					hc.res.Headers["Content-Length"] = fmt.Sprint(len(v))
				}
			}
			hc.Respond()
			return
		}

		pattern := regexp.MustCompile(`/echo/([^/]+)`)
		matches := pattern.FindStringSubmatch(path)
		if len(matches) == 2 {
			hc.res.Body = matches[1]
			hc.res.Headers["Content-Type"] = "text/plain"
			hc.res.Headers["Content-Length"] = fmt.Sprint(len(matches[1]))
			hc.Respond()
			return
		}

		pattern = regexp.MustCompile(`/files/([^/]+)`)
		matches = pattern.FindStringSubmatch(path)
		if len(matches) == 2 {
			filePath := hc.config.DirectoryFlag + "/" + matches[1]
			fileContent, err := readFile(filePath)
			if err != nil {
				hc.res.Status = STATUS_NOT_FOUND
				hc.Respond()
				return
			}

			hc.res.Body = fileContent

			// set headers and respond
			hc.res.Headers["Content-Type"] = "application/octet-stream"
			hc.res.Headers["Content-Length"] = fmt.Sprint(len(fileContent))
			hc.Respond()
			return
		}

	}

	if method == "POST" {
		pattern := regexp.MustCompile(`/files/([^/]+)`)
		matches := pattern.FindStringSubmatch(path)
		if len(matches) == 2 {
			filePath := hc.config.DirectoryFlag + "/" + matches[1]
			fileContent := hc.req.Body

			err := writeFile(filePath, fileContent)
			if err != nil {
				hc.res.Status = STATUS_INTERNAL_SERVER_ERROR
				hc.Respond()
				return
			}

			// set headers and respond
			hc.res.Status = STATUS_CREATED
			hc.Respond()
			return
		}
	}

	hc.res.Status = STATUS_NOT_FOUND
	hc.Respond()
}

// TODO: Research about this more
func readFile(filePath string) (string, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	data := buffer[:n]

	return string(data), nil
}

func writeFile(filePath, content string) error {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
