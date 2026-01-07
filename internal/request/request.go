package request

import (
	"io"
	"strings"
	"fmt"
	"unicode"
)

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

type Request struct {
	RequestLine RequestLine	
	state parserState
}

type parserState int 
const (
	StateInit parserState = 0 
	StateDone parserState = 1
)

func StringIsUpper(s string) bool {
	for _, char := range s {
		if !unicode.IsUpper(char) ||  !unicode.IsLetter(char){
			return false
		}
	}
	return true
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer :
	for { 
		switch r.state {
			case StateInit:
				return read, nil
			case StateDone:
				break outer
		}
	}
	return 0, nil

}

func parseRequestLine(s string) (*RequestLine, int, error) {
	spIndex := strings.Index(s, "\r\n")
	
	if spIndex == -1 {
		return nil, 0, fmt.Errorf("Incomplete request line")
	}
	
	startLine := s[:spIndex]

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
func RequestFromReader(reader io.Reader) (*Request, error) {
	returnedRequest := &Request{ state: StateInit }

	buffer := make([]byte, 1024)
	bIndex := 0

	for {
		reader.Read(buffer[bIndex:])

		n, err := io.Copy(reader)
		bIndex += n
		
		if err != nil {
			return nil, fmt.Errorf("Error in trying to read in the incoming data")
		}

		if returnedRequest.status == StateDone {
			break
		}
	}
	return &Request{
		RequestLine: *returnedRequestLine,
	}, parsingError

}
