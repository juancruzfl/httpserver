package request

import ( 
	"testing"
	"io"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data string
	chunkSize int
	pos int 
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	
	if cr.pos >= len(cr.data) { 
		return 0, io.EOF
	}

	maxBytesRead := cr.chunkSize
	if maxBytesRead > len(p) {
		maxBytesRead = len(p)
	}

	remainingData := len(cr.data) - cr.pos
	if maxBytesRead > remainingData {
		maxBytesRead = remainingData
	}

	n = copy(p, cr.data[cr.pos : cr.pos + maxBytesRead])	
	cr.pos += n

	return n, nil
}

func TestParseRequestLine_NoHeaders(t *testing.T) {
	reader := &chunkReader {
		data: "GET / HTTP/1.1\r\n\r\n",
		chunkSize: 3,
	}

	r, error := RequestFromReader(reader)
	require.NoError(t, error)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
}

func TestParseRequestLine_WithQueryParams(t *testing.T) {
	reader := &chunkReader {
		data: "GET /search?q=test&lang=en HTTP/1.1\r\n\r\n",
		chunkSize: 1,
	}

	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "/search?q=test&lang=en", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestParseRequestLine_HTTP10(t *testing.T) {
	reader := &chunkReader {
		data: "GET /legacy HTTP/1.0\r\n\r\n",
		chunkSize: 4,
	}

	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "/legacy", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.0", r.RequestLine.HttpVersion)
}

func TestParseRequestLine_ExtraSpaces(t *testing.T) {
	reader := &chunkReader {
		data: "GET    /foo/bar    HTTP/1.1\r\n\r\n",
		chunkSize: 3,
	}

	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "/foo/bar", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestParseRequestLine_MissingHTTPVersion(t *testing.T) {
	reader := &chunkReader {
		data: "GET /\r\n\r\n",
		chunkSize: 3,
	}

	r, err := RequestFromReader(reader)
	require.Error(t, err)
	assert.Nil(t, r)
}

func TestParseRequestLine_InvalidHTTPVersion(t *testing.T) {
	reader := &chunkReader {
		data: "GET / HTTP/2\r\n\r\n",
		chunkSize: 3,
	}

	r, err := RequestFromReader(reader)
	require.Error(t, err)
	assert.Nil(t, r)
}

func TestParseRequestLine_EmptyInput(t *testing.T) {
	reader := &chunkReader {
		data: "",
		chunkSize: 3,
	}

	r, err := RequestFromReader(reader)
	require.Error(t, err)
	assert.Nil(t, r)
}

func TestParseRequestLine_TooManyTokens(t *testing.T) {
	reader := &chunkReader {
		data: "GET / HTTP/1.1 EXTRA\r\n\r\n",
		chunkSize: 3,
	}
	r, err := RequestFromReader(reader)
	require.Error(t, err)
	assert.Nil(t, r)

}

func TestParseHeaders(t *testing.T) {
	reader := &chunkReader {
		data: "GET / HTTP/1.1\r\nHost: localhost:8000 \r\nUser-Agent: curl/0.0.0 \r\nAccept: */* \r\nTransfer-Encoding: chunked\r\n\r\n",
		chunkSize: 3,
	}
	r, err := RequestFromReader(reader)
	hostHeader, okhost := r.Headers.Get("Host")
	useragentHeader, okagent := r.Headers.Get("user-agent")
	
	require.NoError(t, err)
	assert.Equal(t, "localhost:8000", hostHeader)
	assert.Equal(t, "curl/0.0.0", useragentHeader)
	assert.True(t, okhost)
	assert.True(t, okagent)
}

func TestValidBodyLength(t *testing.T) {
	reader := &chunkReader {
		data : "GET /test HTTP/1.1\r\nHost: localhost:8091 \r\nUser-Agent: curl/0.0.0 \r\nAccept: image/jpeg \r\nContent-Length: 10\r\n\r\ntesting te",
		chunkSize: 3,
	}
	r, err := RequestFromReader(reader)
	stringBody := string(r.Body)
	contentLength, okcl := r.Headers.Get("Content-Length")

	require.NoError(t, err)
	assert.Equal(t, "10", contentLength)
	assert.True(t, okcl)
	assert.Equal(t, "testing te", stringBody)
}

func TestIncorrectBodyLength(t *testing.T) {
	reader := &chunkReader {
		data : "GET /parsing HTTP/1.1\r\nHost: localhost:8001 \r\nUser-Agent: curl/0.0.0 \r\nAccept: text/html \r\nContent-Length: 9\r\n\r\nsomerandom",
		chunkSize: 3,
	}
	_, err := RequestFromReader(reader)

	require.Error(t, err)
}

func TestBodyOneByteAtATime(t *testing.T) {
    reader := &chunkReader{
        data:      "POST /tiny HTTP/1.1\r\nHost: localhost\r\nContent-Length: 5\r\n\r\nhello",
        chunkSize: 1,
    }
    r, err := RequestFromReader(reader)

    require.NoError(t, err)
    assert.Equal(t, "hello", string(r.Body))
    assert.Equal(t, 5, len(r.Body))
}

func TestZeroLengthBody(t *testing.T) {
    reader := &chunkReader{
        data:      "GET /empty HTTP/1.1\r\nHost: localhost\r\nContent-Length: 0\r\n\r\n",
        chunkSize: 3,
    }
    r, err := RequestFromReader(reader)

    require.NoError(t, err)
    assert.Equal(t, 0, len(r.Body))
    assert.Equal(t, StateDone, r.state)
}

func TestHeaderTooLong(t *testing.T) {
	longHeaderLine := "X-Long: " + strings.Repeat("a", 1500) + "\r\n"
    reader := &chunkReader{
        data:      "GET / HTTP/1.1\r\n" + longHeaderLine + "\r\n",
        chunkSize: 128,
    }
    request, err := RequestFromReader(reader)

    require.Error(t, err)
	assert.Nil(t, request)
}

func TestDuplicateHeaders(t *testing.T) {
    reader := &chunkReader{
        data:      "GET / HTTP/1.1\r\nHost: localhost\r\nSet-Cookie: ID=1\r\nSet-Cookie: User=Admin\r\n\r\n",
        chunkSize: 10,
    }
    r, err := RequestFromReader(reader)

    require.NoError(t, err)
    cookie, ok := r.Headers.Get("Set-Cookie")
    assert.True(t, ok)
    assert.Equal(t, "ID=1, User=Admin", cookie)
}
