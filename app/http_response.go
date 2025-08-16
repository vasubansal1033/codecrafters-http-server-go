package main

import (
	"strconv"
	"strings"
)

const (
	OkResponse       = "HTTP/1.1 200 OK \r\n"
	NotFoundResponse = "HTTP/1.1 404 Not Found \r\n"

	HTTP_VERSION   = "HTTP/1.1"
	CRLF           = "\r\n"
	CONTENT_TYPE   = "Content-Type"
	CONTENT_LENGTH = "Content-Length"
)

type httpResponse struct {
	StatusCode     int
	Headers        map[string]string
	Body           string
	responseString strings.Builder
}

func (response *httpResponse) addHeader(headerName string, headerValue string) {
	if response.Headers == nil {
		response.Headers = map[string]string{}
	}

	response.Headers[headerName] = headerValue
}

func (response *httpResponse) addBody(contentType string, body string) {
	response.addHeader(CONTENT_TYPE, contentType)
	response.addHeader(CONTENT_LENGTH, strconv.Itoa(len(body)))

	response.Body = body
}

func (response *httpResponse) writeHeaders() {
	for key, val := range response.Headers {
		response.responseString.WriteString(key)
		response.responseString.WriteString(": ")
		response.responseString.WriteString(val)
		response.responseString.WriteString(CRLF)
	}

	response.responseString.WriteString(CRLF)
}

func (response *httpResponse) writeStatus() {
	response.responseString.WriteString(HTTP_VERSION)
	response.responseString.WriteString(" ")
	response.responseString.WriteString(strconv.Itoa(response.StatusCode))

	status := ""
	switch response.StatusCode {
	case 200:
		status = "OK"
	case 201:
		status = "Created"
	case 404:
		status = "Not Found"
	}

	if status != "" {
		response.responseString.WriteString(" ")
		response.responseString.WriteString(status)
	}

	response.responseString.WriteString(CRLF)
}

func (response *httpResponse) writeBody() {
	response.responseString.WriteString(response.Body)
}

func (response *httpResponse) toString() (string, error) {
	response.writeStatus()
	response.writeHeaders()
	response.writeBody()

	return response.responseString.String(), nil
}
