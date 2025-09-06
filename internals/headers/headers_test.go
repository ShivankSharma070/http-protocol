package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParsing(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nfoofoo:        barbar  \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	value, _ := headers.Get("host")
	assert.Equal(t, "localhost:42069", value)
	value, _ = headers.Get("foofoo")
	assert.Equal(t, "barbar", value)
	_, ok := headers.Get("MisingKey")
	assert.False(t, ok)
	assert.Equal(t, 50, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid header name
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple values for a single header
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nhost: address.com\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	value, _ = headers.Get("host")
	assert.Equal(t, "localhost:42069, address.com", value)
	assert.Equal(t, 44, n)
	assert.True(t, done)

}
