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

func (h *Headers) Get(name string) (string, bool) {
	key := strings.ToLower(name)
	i, ok := h.headers[key]
	return i, ok
}

func (h *Headers) Set(name string, value string) {
	key := strings.ToLower(name)

	exisiting, ok := h.Get(key)
	
	if ok {
		h.headers[key] = exisiting + ", " + value
	} else {
		h.headers[key] = value
	}
}

func isTokenChar(char byte) bool {
	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
		return true
	}
	
	switch char {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
		return true
	}
	return false 
}

func validateFieldName(name []byte) bool {
	if name == nil {
		return false
	}

	for i := 0; i < len(name); i++ {
		char := name[i]
		if !isTokenChar(char) {
			return false
		}
	}
	return true
}

func parseheaderline(fieldline []byte) (string, string, error) {
	colonIndex := bytes.Index(fieldline, []byte(":"))
	
	if colonIndex == -1 {
		return "", "", fmt.Errorf("Error in trying to parse header line")
	}
	
	lastChar := fieldline[colonIndex - 1]
	if lastChar == ' ' || lastChar == '\t' {
		return "", "", fmt.Errorf("Error in trying to parse header line. Field name contains a tab or space between the colon")
	}

	fieldname := fieldline[:colonIndex]
	fieldvalue := bytes.TrimSpace(fieldline[colonIndex + 1:])

	if valid := validateFieldName(fieldname); !valid{
		return "", "", fmt.Errorf("Error in trying to parse header line. Invalid tokens in field name")
	}

	return string(fieldname), string(fieldvalue), nil
}

func (h *Headers) validateMandatoryHeaders() error {
	_, okContentLength := h.Get("Content-Length")
	_, okTransEncoding := h.Get("Transfer-Encoding")
	exitisingHost, okHost := h.Get("Host")

	if okContentLength && okTransEncoding {
		return fmt.Errorf("Transfer Encoding and Content-Length detected. Not allowed")
	}

	if !okHost {
		return fmt.Errorf("No host detected")
	}
	if strings.Contains(exitisingHost, ",") {
		return fmt.Errorf("More than one host has been detected")
	}
	
	return nil
}


func (h *Headers) Validate() error {
	if err := h.validateMandatoryHeaders(); err != nil {
		return err
	}
	return nil
}

func (h *Headers) parseFieldValues() {
	
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
