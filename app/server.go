package main

import (
	"fmt"
	"net"
	"os"
	"log"
)

const (
	TCP_HOST = "0.0.0.0"
	TCP_PORT = 4221
)

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", TCP_HOST, TCP_PORT))
	if err != nil {
		logAndThrowError(err, fmt.Sprintf("Failed to bind to port: %d", TCP_PORT))
	}
	
	conn, err := l.Accept()
	if err != nil {
		logAndThrowError(err, "Error accepting connection")
	}

	log.Println("Connection accepted")

	statusLine := "HTTP/1.1 200 OK\r\n"
	headers := "\r\n"

	_, err = conn.Write([]byte(statusLine + headers))
	if err != nil {
		logAndThrowError(err, "Error while writing data")
	}
}

func logAndThrowError(err error, errorMessage string) {
	log.Fatalf("%s: %v", errorMessage, err)
	os.Exit(1)
}
