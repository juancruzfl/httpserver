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

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:8000", headers["Host"])
	assert.Equal(t, 23, n)
	assert.True(t, done)
}

func TestInvalidHeaderSpacing(t *testing.T) {
	headers := NewHeaders()
	
	data := []byte("		Host : localhost:8003		\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
