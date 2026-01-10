package headers

import (
	"bytes"
	"strings"
	"fmt"
)

type Headers map[string]string

func NewHeaders() Headers {
	n := make(Headers)
	return n
}

func parseHeaderLine(fieldline []byte) (string, string, error) {
	stringfields := string(fieldline)
	fields := strings.Fields(stringfields)

	for _, field := range fields {
		if field == ":" {
			return "", "", fmt.Errorf("Invalid field line format")
		}
	}

	splitFields := strings.Split(fields[0], ":")
	fieldName := splitFields[0]
	fieldValue := fields[1]

	return fieldName, fieldValue, nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	bytesRead := 0
	done := false

	for {
		sIndex := bytes.Index(data[bytesRead:], []byte("\r\n"))

		if sIndex == -1 {
			break
		}

		if sIndex == 0 {
			done = true
			break
		}

		fieldName, fieldValue, err := parseHeaderLine(data[:sIndex])
	
		if err != nil {
			return 0, done, err
		}
		
		h[fieldName] = fieldValue
		bytesRead += sIndex + 2
	}
	return bytesRead, done, nil
}
