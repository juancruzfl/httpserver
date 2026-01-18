package headers

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidSpacing (t *testing.T) {
	headers :=  NewHeaders()

	data := []byte ("Host: localhost:8000 \r\n\r\n")
	n, done, err := headers.Parse(data)

	headerval, ok := headers.Get("Host")
	mandatoryHeaders := headers.Validate()

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:8000", headerval)
	assert.Equal(t, true, ok)
	assert.Equal(t, 23, n)
	assert.True(t, done)
	assert.NoError(t, mandatoryHeaders)
}

func TestInvalidHeaderSpacing(t *testing.T) {
	headers := NewHeaders()
	
	data := []byte("		Host : localhost:8003		\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestParseValue(t *testing.T) {
	headers := NewHeaders()

	data := []byte ("Host:                developer.mozilla.org             \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 57, n)
	assert.True(t, done)
}

func TestMulipleLines(t *testing.T) {
	headers := NewHeaders()
	
	data := []byte ("Host: localhost:8000 \r\nUser-Agent: curl/0.0.0 \r\nAccept: */* \r\nTransfer-Encoding: chunked\r\n\r\n")
	n, done, err := headers.Parse(data)
	headersvalhost, okhost := headers.Get("host")
	headersvalagent, okagent := headers.Get("User-Agent")
	mandatoryHeaders := headers.Validate()

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.True(t, done)
	assert.Equal(t, true, okhost)
	assert.Equal(t, 90, n)
	assert.Equal(t, "localhost:8000", headersvalhost)
	assert.Equal(t, "curl/0.0.0", headersvalagent)
	assert.Equal(t, true, okagent)
	assert.NoError(t, mandatoryHeaders)
}


func TestMutlipleHeaders(t *testing.T) {
	headers := NewHeaders()
	
	data := []byte("Host: localhost:8000 \r\nAccept: text/html \r\nAccept: application/xhtml+xml \r\nAccept: */*\r\n\r\n")
	n, done, err := headers.Parse(data)
	headersaccept, okaccept := headers.Get("Accept")

	assert.Equal(t, 88, n)
	require.True(t, done)
	require.NoError(t, err)
	require.True(t, okaccept)
	require.Equal(t, "text/html, application/xhtml+xml, */*", headersaccept)
}
