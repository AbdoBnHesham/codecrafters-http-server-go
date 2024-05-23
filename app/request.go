package main

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	Headers map[string]string
	Params  map[string]string
	Method  string
	Path    string
	Body    string
}

func ParseRequest(reader io.Reader) (Request, error) {
	buffer := make([]byte, 1024)
	_, err := reader.Read(buffer)
	if err != nil {
		return Request{}, fmt.Errorf("error reading response: %v", err.Error())
	}

	request := string(buffer)
	sections := strings.Split(request, CRLF+CRLF)
	reqAndHeaders := strings.Split(sections[0], CRLF)

	// parse status line
	requestLineParts := strings.Split(reqAndHeaders[0], " ")
	if len(requestLineParts) < 2 {
		return Request{}, fmt.Errorf("error parsing status line: %s", requestLineParts)
	}

	// parse headers
	headers := make(map[string]string)
	for i := 1; i < len(reqAndHeaders); i++ {
		h := strings.SplitN(reqAndHeaders[i], ":", 2)
		if len(h) == 2 {
			headers[h[0]] = strings.TrimSpace(h[1])
		}
	}

	// TODO parse body

	return Request{
		Method:  requestLineParts[0],
		Path:    requestLineParts[1],
		Headers: headers,
		Body:    sections[1],
	}, nil
}
