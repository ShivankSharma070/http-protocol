package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/ShivankSharma070/http-protocol/internals/headers"
)

type parserState string

const (
	StateInit   parserState = "init"
	StateDone   parserState = "done"
	StateHeader parserState = "header"
	StateError  parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (re *RequestLine) verifyRequestLine() error {
	if !(re.Method == "GET" || re.Method == "POST" || re.Method == "DELETE" || re.Method == "PUT") {
		return ERROR_UNSUPPORTED_METHOD
	}

	httpParts := strings.Split(re.HttpVersion, "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return ERROR_UNSUPPORTED_HTTP_VERSION
	}

	re.HttpVersion = httpParts[1]
	return nil
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	State       parserState
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		switch r.State {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			re, n, err := parseRequestLine(currentData)
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *re
			read += n

			r.State = StateHeader

		case StateHeader:
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				r.State = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if done {
				r.State = StateDone
			}

		case StateDone:
			break outer

		default:
			panic("Somehow we programmed very badly")
		}
	}
	return read, nil
}

func (re *Request) done() bool {
	return re.State == StateDone
}

func newRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("Malformed Request-Line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("Unsupported http version")
var ERROR_UNSUPPORTED_METHOD = fmt.Errorf("Unsupported method type")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("Request in Error state")
var SEPERATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startMsg := b[:idx]
	noOfBytesParsed := idx + len(SEPERATOR)
	parts := bytes.Split(startMsg, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	re := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(parts[2]),
	}

	if err := re.verifyRequestLine(); err != nil {
		return nil, 0, err
	}

	return re, noOfBytesParsed, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	// NOTE: Buffer could get overrun... a header that's above 1k could do that or the body
	buf := make([]byte, 1024)
	bufLen := 0

	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		// TODO: What to do here?
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])

		if err != nil {
			return nil, err
		}
		// Remove the already parsed data from buf
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
