package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	TCP_HOST = "0.0.0.0"
	TCP_PORT = 4221

	ECHO_PATH       = "/echo/"
	USER_AGENT_PATH = "/user-agent"
	FILE_PATH       = "/files/"

	USER_AGENT       = "User-Agent"
	PLAIN_TEXT       = "text/plain"
	OCTET_STREAM     = "application/octet-stream"
	ACCEPT_ENCODING  = "Accept-Encoding"
	CONTENT_ENCODING = "Content-Encoding"
	GZIP             = "gzip"
)

var WORKING_DIRECTORY string

// read one request at a time from connection
func readRequestString(conn net.Conn) string {
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	reader := bufio.NewReader(conn)

	// Read the request line first
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return ""
		}
		logAndThrowError(err, "Failed to read request line")
	}

	requestString := strings.TrimRight(requestLine, CRLF)

	// read headers
	contentLength := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			logAndThrowError(err, "Failed to read request bytes")
		}

		requestString += line

		// Check if this line is a header or empty line
		if line == CRLF {
			break
		}

		if strings.HasPrefix(strings.ToLower(line), strings.ToLower(CONTENT_LENGTH)) {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				contentLength, err = strconv.Atoi(strings.TrimSpace(parts[1]))
				if err != nil {
					logAndThrowError(err, "Failed to parse content length")
				}
			}
		}
	}

	// Read body if Content-Length is specified
	if contentLength > 0 {
		body := make([]byte, contentLength)
		_, err := io.ReadFull(reader, body)
		if err != nil {
			logAndThrowError(err, "Failed to read body")
		}

		requestString += string(body)
	}

	return requestString
}

// TODO: refactor to map a handler func to each request path and verb [like net/http]
func respondToHttpRequest(conn net.Conn, r *httpRequest) {
	response := &httpResponse{}

	// Add Connection header based on request
	if shouldCloseConnection(r) {
		response.addHeader("Connection", "close")
	} else {
		response.addHeader("Connection", "keep-alive")
	}

	if r.Path == "/" {
		response.StatusCode = 200
	} else if strings.HasPrefix(r.Path, ECHO_PATH) {
		response.StatusCode = 200
		responseBody := getEchoResponseBody(r, response)

		response.addBody(PLAIN_TEXT, responseBody)

	} else if strings.HasPrefix(r.Path, USER_AGENT_PATH) {
		response.StatusCode = 200
		responseBody := r.Headers[USER_AGENT]

		response.addBody(PLAIN_TEXT, responseBody)
	} else if strings.HasPrefix(r.Path, FILE_PATH) {
		filePath := path.Join(WORKING_DIRECTORY, r.Path[len(FILE_PATH):])

		if r.Method == "GET" {
			if _, err := os.Stat(filePath); err != nil {
				response.StatusCode = 404
			} else {
				response.StatusCode = 200
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					logAndThrowError(err, fmt.Sprintf("Error while reading file: %s", filePath))
				}

				responseBody := string(fileContent)
				response.addBody(OCTET_STREAM, responseBody)
			}
		} else if r.Method == "POST" {
			response.StatusCode = 201

			f, err := os.Create(filePath)
			defer func() {
				err := f.Close()
				if err != nil {
					logAndThrowError(err, "Error while closing the file")
				}
			}()

			if err != nil {
				logAndThrowError(err, fmt.Sprintf("Error while creating file: %s", filePath))
			}

			w := bufio.NewWriter(f)
			w.Write(r.Body)
			w.Flush()
		}
	} else {
		// not found
		response.StatusCode = 404
	}

	responseString, err := response.toString()
	if err != nil {
		logAndThrowError(err, "Error while creating response string.")
	}

	_, err = conn.Write([]byte(responseString))
	if err != nil {
		logAndThrowError(err, "Error while writing response data")
	}
}

func getEchoResponseBody(r *httpRequest, response *httpResponse) string {
	if r.Method == "GET" {
		if _, ok := r.Headers[ACCEPT_ENCODING]; ok {
			encodingHeaders := strings.Split(r.Headers[ACCEPT_ENCODING], ",")
			for _, encodingHeader := range encodingHeaders {
				if strings.TrimSpace(encodingHeader) == GZIP {
					response.addHeader(CONTENT_ENCODING, GZIP)

					responseBody := r.Path[len(ECHO_PATH):]

					var compressedBody bytes.Buffer
					gz := gzip.NewWriter(&compressedBody)
					_, err := gz.Write([]byte(responseBody))
					if err != nil {
						logAndThrowError(err, "Error while writing to gzip writer")
					}

					err = gz.Close()
					if err != nil {
						logAndThrowError(err, "Error while closing gzip writer")
					}

					return compressedBody.String()
				}
			}
		}
	}

	responseBody := r.Path[len(ECHO_PATH):]
	return responseBody
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		requestString := readRequestString(conn)
		if requestString == "" {
			// connection closed by client
			return
		}

		httpRequest := newHttpRequest(requestString)
		shouldClose := shouldCloseConnection(httpRequest)

		respondToHttpRequest(conn, httpRequest)

		if shouldClose {
			return
		}
	}
}

func shouldCloseConnection(r *httpRequest) bool {
	// Check Connection header
	if connectionHeader, exists := r.Headers["Connection"]; exists {
		return strings.ToLower(strings.TrimSpace(connectionHeader)) == "close"
	}

	// Check HTTP version - HTTP/1.0 defaults to close, HTTP/1.1 defaults to keep-alive
	if r.HttpVersion == "HTTP/1.0" {
		return true
	}

	// HTTP/1.1 defaults to keep-alive unless Connection: close is specified
	return false
}

func main() {
	if len(os.Args) < 2 {
		WORKING_DIRECTORY, _ = os.Getwd()
	} else {
		WORKING_DIRECTORY = os.Args[2]
	}

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", TCP_HOST, TCP_PORT))
	if err != nil {
		logAndThrowError(err, fmt.Sprintf("Failed to bind to port: %d", TCP_PORT))
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logAndThrowError(err, "Error accepting connection")
		}

		log.Println("Connection accepted")

		// handle each incoming connection in new goroutine
		go handleConnection(conn)
	}
}

func logAndThrowError(err error, errorMessage string) {
	log.Fatalf("%s: %v", errorMessage, err)
	os.Exit(1)
}
