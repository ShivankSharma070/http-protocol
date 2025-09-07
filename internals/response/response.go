package response

import (
	"fmt"
	"io"

	"github.com/ShivankSharma070/http-protocol/internals/headers"
)

type Response struct {
}

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-type", "text/plain")
	return h
}


type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer{
	return &Writer{
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOK:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("unrecognized status code")
	}
	_, err := w.writer.Write(statusLine)

	return err
}

func (w *Writer)  WriteHeaders(h *headers.Headers) error {
	b := []byte{}
	var err error = nil
	h.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})

	b = fmt.Appendf(b, "\r\n")
	_, err = w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(b []byte) (int, error) {
	return w.writer.Write(b)
}
