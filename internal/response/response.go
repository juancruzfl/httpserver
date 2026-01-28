package response

import (
	"github.com/juancruzfl/httpserver/internal/headers"
)

type CustomeResponseWriter interface {
	Headers headers.Headers
	Write([]byte) (int, error)
	CustomWriteHeader(statusCode int)
}

type Response struct {
	Status int
	Headers headers.Headers
	Body []byte
	HttpVersion string
}

func newResponse() *Response {
	return &Response{
		Status: 200,
		Heeders: *headers.NewHeaders(),
		Body: nil,
		HttpVersion: "",
	}
}
