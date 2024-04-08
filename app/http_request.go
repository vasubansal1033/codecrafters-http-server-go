package main

import (
	"bufio"
	"strings"
)

type httpRequest struct {
	Method      string
	Path        string
	HttpVersion string
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

func parseHttpRequest(requestString string) *httpRequest {
	httpRequest := &httpRequest{}

	// buffered way of reading the request string line by line
	requestScanner := bufio.NewScanner(strings.NewReader(requestString))
	requestScanner.Split(bufio.ScanLines)

	requestScanner.Scan()
	requestStartLine := requestScanner.Text()

	httpRequest.parseStartLine(requestStartLine)
	return httpRequest
}

func newHttpRequest(requestString string) *httpRequest {
	return parseHttpRequest(requestString)
}
