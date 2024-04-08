package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	TCP_HOST = "0.0.0.0"
	TCP_PORT = 4221

	ECHO_PATH       = "/echo/"
	USER_AGENT_PATH = "/user-agent"
	USER_AGENT      = "User-Agent"
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
	response := &httpResponse{}
	if r.Path == "/" {
		response.StatusCode = 200
	} else if strings.HasPrefix(r.Path, ECHO_PATH) {
		response.StatusCode = 200
		responseBody := r.Path[len(ECHO_PATH):]
		response.addBody("text/plain", responseBody)
	} else if strings.HasPrefix(r.Path, USER_AGENT_PATH) {
		response.StatusCode = 200
		responseBody := r.Headers[USER_AGENT]

		response.addBody("text/plain", responseBody)
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

func handleConnection(conn net.Conn) {
	requestString := readRequestString(conn)
	httpRequest := newHttpRequest(requestString)

	respondToHttpRequest(conn, httpRequest)

	conn.Close()
}

func main() {
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
