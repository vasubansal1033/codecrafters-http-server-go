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
	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", TCP_HOST, TCP_PORT))
	if err != nil {
		logAndThrowError(err, fmt.Sprintf("Failed to bind to port: %d", TCP_PORT))
	}
	
	_, err = conn.Accept()
	if err != nil {
		logAndThrowError(err, "Error accepting connection")
	}
}

func logAndThrowError(err error, errorMessage string) {
	log.Fatalf("%s: %v", errorMessage, err)
	os.Exit(1)
}
