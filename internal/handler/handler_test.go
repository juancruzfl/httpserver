package handler

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/juancruzfl/httpserver/internal/request"
	"github.com/juancruzfl/httpserver/internal/headers"
)

type MockResponseWriter struct {
	Status int
	Headers *headers.Headers
	Body []byte
}

func (m *MockResponseWriter) GetHeaders() *headers.headers {
	return m.Headers
}

func (m *MockResponseWriter) ResponseWriter(status int) {
	m.Status = status
}

func (m *MockResponseWriter) Write(data []byte) (int, error) {
	m.Body = append(m.body, data...)
	return len(data), nil
}

func MockResponseWriter() MockResponseWriter {
	return &MockResponseWriter{
	Status: 200, 
	Headers: headers.NewHeaders(),
	}
}

func TestMockHandler(t *testing.T) {
	mockWriter := MockResponseWriter()
	mockReqLine := request.RequestLine{
		HttpVersion: "1.1",
		RequestTarget: "/testing",
		Method: "GET",
	}
	mockReq := request.Request{
		RequestLine: mockReqLine,
		Headers headers.Headers
		Body []byte
		bodyLength int
		state parserState
	}
	var myHandler HandlerFunc = func(w response.ResponseWriter, r *request.Request) {
			w.CustomWriteHeader(200)
			w.Write([]byte("Hello, Neovim!"))
		}
	myHandler.ServeHttp(mockWriter, mockReq)
	if mockWriter.Status != 200 {
			t.Errorf("Expected status 200, got %d", mockWriter.Status)
		}
    if string(mockWriter.Body) != "Hello, Neovim!" {
        t.Errorf("Expected body 'Hello, Neovim!', got %s", string(mockWriter.Body))
    }
}
