package headers

import (
	"bytes"
	"strings"
	"fmt"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name string, value string) {
	h.headers[strings.ToLower(name)] = value
}

func parseheaderline(fieldline []byte) (string, string, error) {
	stringfields := string(fieldline)
	fields := strings.Fields(stringfields)

	for _, field := range fields {
		if field == ":" {
			return "", "", fmt.Errorf("invalid field line format")
		}
	}

	splitfields := strings.Split(fields[0], ":")
	fieldname := splitfields[0]
	fieldvalue := fields[1]

	return fieldname, fieldvalue, nil
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	bytesread := 0
	done := false

	for {
		sindex := bytes.Index(data[bytesread:], []byte("\r\n"))

		if sindex == -1 {
			break
		}

		if sindex == 0 {
			done = true
			break
		}

		fieldname, fieldvalue, err := parseheaderline(data[bytesread : bytesread + sindex])
	
		if err != nil {
			return 0, done, err
		}
		
		h.Set(fieldname, fieldvalue)
		bytesread += sindex + 2
	}
	return bytesread, done, nil
}
