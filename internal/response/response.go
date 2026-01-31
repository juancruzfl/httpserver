package response

import (
	"net"
	"github.com/juancruzfl/httpserver/internal/headers"
	"github.com/juancruzfl/httpserver/internal/request"
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
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, phrase)
	r.conn.Write([]byte(statusLine))
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
    case 500: 
		phrase = "Internal Server Error"
    default:  
		phrase = "OK"
    }
	r.StatusLine.Status = status
	r.StatusLine.StatusPhrase = phrase
	for key, value := range r.Headers.headers {
		headerLine := fmt.Sprintf("%s: %s\r\n", key, value)
		r.conn.Write([]byte(headerLine))
	}
	// this seperator is used to keep our headers seperate from our body when writing the resonse
	r.conn.Write([]byte("\r\n"))
	r.WritingState = true
}

func (r *Response) Write(data []byte) (int, error) {
	if !r.WritingState {
		r.CustomWriteHeader(200)
	}
	return r.conn.Write(data)
}

func NewResponseWriter(conn net.Conn) ResponseWriter {
	return &Response{
		conn: conn,
		StatusLine: StateLine{
			Status: 200,
			HttpVersion: "",
			StatusPhrase: "OK",
		},
		Headers: headers.NewHeaders(),
		HttpVersion: "",
		WritingState: false,
	}
}

