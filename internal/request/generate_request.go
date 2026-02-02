package request

import (
	"github.com/juancruzfl/httpserver/internal/headers"
	"math/rand/v2"
)

var patterns = [12]string{"/testing1", "/testing2", "/testing3", "/users", "/addresses", "/emails", "/devices", "/users/v1", "/email/v2", "/school", "/tokens/v1", "/tokens/v3"}
var methods = [10]string{"GET", "POST", "PUT", "DELETE", "CONNECT", "HEAD", "OPTIONS", "PATCH", "PUT", "TRACE"}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func NewMockRequestLine(version, pattern, method string) *RequestLine {
	return &RequestLine{
		HttpVersion: version,
		RequestTarget: pattern,
		Method: method,
	}
}

func NewRandomRequestLine() *RequestLine{
	randomPatternIndex := randRange(0, 12)
	randomMethodIndex := randRange(0, 10)
	var requestLine = NewMockRequestLine("1.1", patterns[randomPatternIndex], methods[randomMethodIndex])
	return requestLine
}

func NewMockRequest(requestLine *RequestLine,headers headers.Headers, body []byte, bodyLength int) *Request {
	request := newRequest()
	request.RequestLine = *requestLine
	request.Headers = headers
	request.Body = body
	request.bodyLength = bodyLength 
	return request
}

