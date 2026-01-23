package response

import (
	"github.com/juancruzfl/httpserver/internal/headers"
)

type Response struct {
	Status int
	Headers headers.Headers
	Body []byte
	HttpVersion string
}

func newResponse() *Response {
	return &Response{
		Status: 200,
		Headers: *headers.NewHeaders(),
		Body: nil,
		HttpVersion: "",
	}
}
