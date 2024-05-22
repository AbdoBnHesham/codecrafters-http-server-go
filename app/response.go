package main

import "fmt"

const (
	HTTP_VERSION                 = "HTTP/1.1"
	STATUS_OK                    = "200 OK"
	STATUS_NOT_FOUND             = "404 Not Found"
	STATUS_INTERNAL_SERVER_ERROR = "500 Internal Server Error"
	CRLF                         = "\r\n"
)

type Response struct {
	Headers map[string]string
	Status  string
	Body    string
}

func newResponse() Response {
	return Response{Status: STATUS_OK, Headers: make(map[string]string)}
}

func (r Response) ToBytes() []byte {
	statusLine := fmt.Sprintf("%s %s", HTTP_VERSION, r.Status)

	headers := ""
	for k, v := range r.Headers {
		h := fmt.Sprintf("%s%s %s", k, ":", v)
		headers = fmt.Sprintf("%s%s%s", headers, CRLF, h)
	}
	if headers == "" {
		headers = fmt.Sprintf("%s%s", headers, CRLF)
	}
	headers = fmt.Sprintf("%s%s", headers, CRLF)

	body := ""
	if r.Body != "" {
		body = fmt.Sprintf("%s%s", CRLF, r.Body)
	}
	return []byte(fmt.Sprintf("%s%s%s", statusLine, headers, body))
}
