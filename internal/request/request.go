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
}

func StringIsUpper(s string) bool {
	for _, char := range s {
		if !unicode.IsUpper(char) ||  !unicode.IsLetter(char){
			return false
		}
	}
	return true
}

func ParseRequestLine(s string) (string, string, string, error) {
	
	splitString1 := strings.Split(s, " ")

	if len(splitString1) < 2 {
		return "", "", "", fmt.Errorf("Request line is missing information")
	}

	method := splitString1[0]

	if StringIsUpper(method) == false {
		return "", "", "", fmt.Errorf("Error in the request line parsing: Invalid Method format")
	}

	requestTarget := splitString1[1]

	if strings.HasPrefix(requestTarget, "/") == false {
		return "", "", "", fmt.Errorf("Error in the request line parsing: Invalid resource format")
	}

	versionSplit1 := strings.Split(splitString1[2], "/")

	versionSplit2 := strings.Split(versionSplit1[1], "\r\n")

	httpVersion := versionSplit2[0]

	if httpVersion != "1.1" {
		return "", "", "", fmt.Errorf("Error in the request line parsing: Unsupported http version")
	}

	return method, requestTarget, httpVersion, nil

}
func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	requestLineString := string(buffer)

	requestMethod, requestTarget, httpVersion, parsingError := ParseRequestLine(requestLineString)

	if parsingError != nil {
		return nil, parsingError
	}

	var returnedRequest Request

	returnedRequest.RequestLine.HttpVersion = httpVersion
	returnedRequest.RequestLine.RequestTarget = requestTarget
	returnedRequest.RequestLine.Method = requestMethod

	return &returnedRequest, nil
 
}
