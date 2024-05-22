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
}

func ParseRequest(reader io.Reader) (Request, error) {
	buffer := make([]byte, 1024)
	_, err := reader.Read(buffer)
	if err != nil {
		return Request{}, fmt.Errorf("error reading response: %v", err.Error())
	}

	request := string(buffer)
	lines := strings.Split(request, CRLF)

	// parse status line
	requestLineParts := strings.Split(lines[0], " ")
	if len(requestLineParts) < 2 {
		return Request{}, fmt.Errorf("error parsing status line: %s", requestLineParts)
	}

	// TODO parse headers
	// TODO parse body

	return Request{
		Method: requestLineParts[0],
		Path:   requestLineParts[1],
	}, nil
}
