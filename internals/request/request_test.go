package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data           string
	numBytesToRead int
	pos            int
}

// Reads reads upto len(b) or numBytesToRead bytes from the string per call
// It simulates reading a variable lenght of chunk from a network connection
func (cr *chunkReader) Read(b []byte) (int, error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}

	endIndex := min(cr.pos+cr.numBytesToRead, len(cr.data))

	n := copy(b, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func TestRequestLineParsing(t *testing.T) {
	reader := &chunkReader{
		data:           "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesToRead: 3,
	}

	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	value, _ := r.Headers.Get("host")
	assert.Equal(t, "localhost:42069", value)
	value, _ = r.Headers.Get("user-agent")
	assert.Equal(t, "curl/7.81.0", value)
	value, _ = r.Headers.Get("accept")
	assert.Equal(t, "*/*", value)

	// Empty headers
	reader = &chunkReader{
		data:           "GET / HTTP/1.1\r\n\r\n",
		numBytesToRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)

	// Duplicate header
	reader = &chunkReader{
		data:           "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nhost: another\r\n\r\n",
		numBytesToRead: 1,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	value, _ = r.Headers.Get("host")
	assert.Equal(t, "localhost:42069, another", value)

	// Malformed Request-line
	reader = &chunkReader{
		data:           "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesToRead: 2,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)

	// Unkown http method
	reader = &chunkReader{
		data:           "SET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesToRead: 10,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)

	// Malformed header line
	reader = &chunkReader{
		data:           "SET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent    : curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesToRead: 10,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)

	// Mising end of headers
	reader = &chunkReader{
		data:           "SET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n",
		numBytesToRead: 10,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)
}

func TestRequestBody(t *testing.T) {
	// Test: Standard Body
	reader := &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 13\r\n" +
			"\r\n" +
			"hello world!\n",
		numBytesToRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "hello world!\n", r.Body)

	// Test: Body shorter than reported content length
	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 20\r\n" +
			"\r\n" +
			"partial content",
		numBytesToRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)
}
