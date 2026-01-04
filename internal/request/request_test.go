package request

import ( 
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRequestLine_NoHeaders (t *testing.T) {
	r, error := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\n\r\n"))
	require.NoError(t, error)
	require.NotNil(t, r)
	assert.Equal(t, "GE", r.RequestLine.Method)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
}

func TestParseRequestLine_WithQueryParams(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"GET /search?q=test&lang=en HTTP/1.1\r\n\r\n",
	))
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "/search?q=test&lang=en", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestParseRequestLine_HTTP10(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"GET /legacy HTTP/1.0\r\n\r\n",
	))
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "/legacy", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.0", r.RequestLine.HttpVersion)
}

func TestParseRequestLine_ExtraSpaces(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"GET    /foo/bar    HTTP/1.1\r\n\r\n",
	))
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "/foo/bar", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestParseRequestLine_MissingHTTPVersion(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"GET /\r\n\r\n",
	))
	require.Error(t, err)
	assert.Nil(t, r)
}

func TestParseRequestLine_InvalidHTTPVersion(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"GET / HTTP/2\r\n\r\n",
	))
	require.Error(t, err)
	assert.Nil(t, r)
}

func TestParseRequestLine_EmptyInput(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(""))
	require.Error(t, err)
	assert.Nil(t, r)
}

func TestParseRequestLine_TooManyTokens(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"GET / HTTP/1.1 EXTRA\r\n\r\n",
	))
	require.Error(t, err)
	assert.Nil(t, r)
}
