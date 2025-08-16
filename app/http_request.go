package main

import (
	"bufio"
	"bytes"
	"strings"
)

type httpRequest struct {
	Method      string
	Path        string
	HttpVersion string
	Headers     map[string]string

	Body []byte
}

func (r *httpRequest) parseStartLine(requestStartLine string) {
	scanner := bufio.NewScanner(strings.NewReader(requestStartLine))
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	r.Method = scanner.Text()

	scanner.Scan()
	r.Path = scanner.Text()

	scanner.Scan()
	r.HttpVersion = scanner.Text()
}

func (r *httpRequest) parseHeaderLine(headerLine string) {
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}

	parts := strings.Split(headerLine, ":")
	r.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
}

func parseHttpRequest(requestString string) *httpRequest {
	httpRequest := &httpRequest{}

	// buffered way of reading the request string line by line
	requestScanner := bufio.NewScanner(strings.NewReader(requestString))
	requestScanner.Split(bufio.ScanLines)

	requestScanner.Scan()
	requestStartLine := requestScanner.Text()

	httpRequest.parseStartLine(requestStartLine)

	// read header lines
	for requestScanner.Scan() {
		headerLine := requestScanner.Text()
		if headerLine == "" {
			break
		}

		httpRequest.parseHeaderLine(headerLine)
	}

	requestScanner.Scan()
	httpRequest.Body = requestScanner.Bytes()
	httpRequest.Body = bytes.ReplaceAll(httpRequest.Body, []byte{0}, []byte{})

	return httpRequest
}

func newHttpRequest(requestString string) *httpRequest {
	return parseHttpRequest(requestString)
}
