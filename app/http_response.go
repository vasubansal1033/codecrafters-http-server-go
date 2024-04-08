package main

const (
	OkResponse       = "HTTP/1.1 200 OK \r\n"
	NotFoundResponse = "HTTP/1.1 404 Not Found \r\n"
)

type httpResponse struct {
	StatusCode     int
	Headers        map[string]string
	Body           string
	resopnseString string
}
