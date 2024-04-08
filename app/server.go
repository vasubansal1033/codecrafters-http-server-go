package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	// "strings"
)

const (
	TCP_HOST = "0.0.0.0"
	TCP_PORT = 4221
)

func readRequestString(conn net.Conn) string {
	readBuffer := make([]byte, 1024)
	_, err := conn.Read(readBuffer)
	if err != nil {
		logAndThrowError(err, "Failed to read request bytes")
	}

	return string(readBuffer)
}

func respondToHttpRequest(conn net.Conn, r *httpRequest) {
	response := NotFoundResponse
	headers := "\r\n"
	responseBody := ""
	if r.Path == "/" {
		response = OkResponse
	} else if strings.HasPrefix(r.Path, "/echo/") {
		responseContent := r.Path[len("/echo/"):]
		response = OkResponse

		headers = "Content-Type: text/plain\r\n" + fmt.Sprintf("Content-Length: %d\r\n", len(responseContent)) + "\r\n"
		responseBody = responseContent + "\r\n"
	}

	_, err := conn.Write([]byte(response + headers + responseBody))
	if err != nil {
		logAndThrowError(err, "Error while writing data")
	}
}

func handleConnection(conn net.Conn) {
	requestString := readRequestString(conn)
	httpRequest := newHttpRequest(requestString)

	respondToHttpRequest(conn, httpRequest)
}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", TCP_HOST, TCP_PORT))
	if err != nil {
		logAndThrowError(err, fmt.Sprintf("Failed to bind to port: %d", TCP_PORT))
	}

	conn, err := l.Accept()
	if err != nil {
		logAndThrowError(err, "Error accepting connection")
	}
	defer conn.Close()

	log.Println("Connection accepted")
	handleConnection(conn)
}

func logAndThrowError(err error, errorMessage string) {
	log.Fatalf("%s: %v", errorMessage, err)
	os.Exit(1)
}
