package request

import (
	"io"
	"strings"
	"strconv"
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
	Body []byte
	bodyLength int
	state parserState
}

func newRequest() *Request {
	return &Request{
		Headers: *headers.NewHeaders(),
		state: StateInit,
		Body: nil,
		bodyLength: 0,
	}
}
type parserState int 
const (
	StateInit parserState = 0 
	StateDone parserState = 1
	StateHeaders parserState = 2
	StateBodyFixed parserState = 3
	StateBodyChunked parserState = 4
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
		if len(data[bytesRead:]) == 0 {
			break outer
		}
		
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
					if contentLength, clok := r.Headers.Get("Content-Length"); clok {
						bodyLength, bodyErr := strconv.Atoi(contentLength)
						if bodyErr != nil {
							return 0, bodyErr
						}
						if bodyLength == 0 {
							r.state = StateDone
							break outer
						}
						r.bodyLength = bodyLength
						r.state = StateBodyFixed
					} else if _, teok := r.Headers.Get("Transfer-Encoding"); teok {
						r.state = StateDone
					} else {
						r.state = StateDone
					}
				}
				
				return bytesRead, nil
			case StateBodyFixed:
				remaining := r.bodyLength - len(r.Body)
				available := len(data) - bytesRead		
				if available > remaining {
					return bytesRead, fmt.Errorf("Too many characters sent by request")
				}
				
				r.Body = append(r.Body, data[bytesRead: bytesRead + available]...)
				bytesRead += available
				
				if len(r.Body) == r.bodyLength {
					r.state = StateDone
				}	
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
		
		
		if numRead == 0 && bSize == len(buffer) {
			return nil, fmt.Errorf("Request header on line too long")
		}
	}

	return request, nil
}
