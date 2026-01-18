package request

import (
	"io"
	"strings"
	"errors"
	"fmt"
	"unicode"
	"bytes"
	"github.com/juancruzfl/httpserver/internal/headers"
)

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

type Request struct {
	RequestLine RequestLine	
	Headers headers.Headers
	state parserState
}

func newRequest() *Request {
	return &Request{
		Headers: *headers.NewHeaders(),
		state: StateInit,
	}
}
type parserState int 
const (
	StateInit parserState = 0 
	StateDone parserState = 1
	StateHeaders parserState = 2
)

func StringIsUpper(s string) bool {
	for _, char := range s {
		if !unicode.IsUpper(char) ||  !unicode.IsLetter(char){
			return false
		}
	}
	return true
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	spIndex := bytes.Index(b, []byte("\r\n"))

	if spIndex == -1 {
		return nil, 0, nil
	}

	startLine := string(b[:spIndex])

	requestLineFields := strings.Fields(startLine)

	if len(requestLineFields) != 3 {
		return nil, 0, fmt.Errorf("Request line is missing information")
	}

	method := requestLineFields[0]

	if StringIsUpper(method) == false {
		return nil, 0, fmt.Errorf("Error in the request line parsing: Invalid Method format")
	}

	requestTarget := requestLineFields[1]

	if strings.HasPrefix(requestTarget, "/") == false {
		return nil, 0, fmt.Errorf("Error in the request line parsing: Invalid resource format")
	}

	versionSplit1 := strings.Split(requestLineFields[2], "/")

	versionSplit2 := strings.Split(versionSplit1[1], "\r\n")

	httpVersion := versionSplit2[0]

	if httpVersion != "1.1" && httpVersion != "1.0" && httpVersion != "2.0" && httpVersion != "3.0" {
		return nil, 0, fmt.Errorf("Error in the request line parsing: Unsupported http version")
	}

	var returnedRequestLine RequestLine 

	returnedRequestLine.HttpVersion = httpVersion
	returnedRequestLine.RequestTarget = requestTarget
	returnedRequestLine.Method = method

	return &returnedRequestLine, spIndex, nil

}

func (r *Request) parse(data []byte) (int, error) {
	bytesRead := 0
	outer :
		for { 
			switch r.state {
				case StateInit:
					requestLine, n, err := parseRequestLine(data[bytesRead:])
					if err != nil {
						return 0, err
					}
					
					if n == 0 {
						break outer
					}

					r.RequestLine = *requestLine
					bytesRead += n + 2
					r.state = StateHeaders
				
				case StateHeaders:
					headerBytesRead, done, err := r.Headers.Parse(data[bytesRead:])
					if err != nil {
						return 0, err
					}

					bytesRead += headerBytesRead
					
					if done {
						r.state = StateDone
					}
					return bytesRead, nil
				case StateDone:
					break outer
			}
		}
	return bytesRead, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buffer := make([]byte, 1024)
	bSize := 0

	for request.state != StateDone {
		n, err := reader.Read(buffer[bSize:])
		if err != nil {
			if errors.Is(err, io.EOF){
				if request.state == StateDone {	
					return request, nil
				}
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
		bSize += n

		numRead, error :=  request.parse(buffer[:bSize])
		if error != nil {
			return nil, error
		}

		nCopy := copy(buffer, buffer[numRead: bSize])
		bSize = nCopy

	}

	return request, nil
}
