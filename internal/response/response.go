package response

import (
	"fmt"
	"net"
	"github.com/juancruzfl/httpserver/internal/headers"
)
// sidenote: In the offical documentation for the net/http package, GetHeaders is called Headers and it returns the map of Headers.
// However, in my implementation we will just return a pointer to the Header struct and then use our methods to access the map.
type ResponseWriter interface {
	GetHeaders() *headers.Headers
	Write([]byte) (int, error)
	CustomWriteHeader(int)
}

type StatusLine struct {
	Status int
	HttpVersion string
	StatusPhrase string
}

type Response struct {
	conn net.Conn
	StatusLine StatusLine
	Headers *headers.Headers
	Body []byte
	WritingState bool
}

// this is a similar to our Get function, except we return the header struct, again as a pointer to avoid large copies, instead of an individual one
func (r *Response) GetHeaders() *headers.Headers {
	return r.Headers
}

// sidenote: I am not implementing the full range of status codes. There are too many to do in a short amount of time I have but the essential status codes that are needed for the response to fail or continue
// will take priority.
func (r *Response) CustomWriteHeader(status int) {
	if r.WritingState {
		return 
	}
	var phrase string
    switch status {
    case 200:
		phrase = "OK"
    case 201: 
		phrase = "Created"
    case 400: 
		phrase = "Bad Request"
    case 404: 
		phrase = "Not Found"
    case 501: 
		phrase = "Internal Server Error"
    default:  
		phrase = "OK"
    }
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, phrase)
	r.conn.Write([]byte(statusLine))
	r.StatusLine.Status = status
	r.StatusLine.StatusPhrase = phrase

	requestHeaders := r.Headers
	
	requestHeaders.ForEach(func(key, value string) {
		headerLine := fmt.Sprintf("%s: %s\r\n", key, value)
		r.conn.Write([]byte(headerLine))
	})
	// this seperator is used to keep our headers seperate from our body when writing the resonse
	r.conn.Write([]byte("\r\n"))
	r.WritingState = true
}

func (r *Response) Write(data []byte) (int, error) {
	if !r.WritingState {
		r.CustomWriteHeader(200)
	}
	n, err := r.conn.Write(data)
	if err != nil {
		fmt.Errorf("Error while to trying to write to response: %w ", err)
	}
	return n, nil
}


// huge sidenote: I am taking the approach of using a struct that implements some of the methods of a the ResponseWriter interface. Some are decprecated, others are not. This may/will introduce bugs in more complex http server implementations.
// If you are interestead in reading more about this, check this out: https://github.com/felixge/httpsnoop?tab=readme-ov-file#why-this-package-exists
func NewResponseWriter(conn net.Conn) ResponseWriter {
	return &Response{
		conn: conn,
		StatusLine: StatusLine{
			Status: 200,
			HttpVersion: "",
			StatusPhrase: "OK",
		},
		Headers: headers.NewHeaders(),
		WritingState: false,
	}
}

