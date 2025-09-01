package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
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
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("Malformed Request-Line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("Unsupported http version")
var ERROR_UNSUPPORTED_METHOD = fmt.Errorf("Unsupported method type")
var SEPERATOR = "\r\n"

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, b, nil
	}

	startMsg := b[:idx]
	restMsg := b[idx+len(SEPERATOR):]
	parts := strings.Split(startMsg, " ")

	if len(parts) != 3 {
		return nil, restMsg, ERROR_MALFORMED_REQUEST_LINE
	}

	re := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}

	if err := re.verifyRequestLine(); err != nil {
		return nil, restMsg, err
	}

	return re, restMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("unable to io.ReadAll"),
			err,
		)
	}

	rq, _, err := parseRequestLine(string(data))
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rq,
	}, nil
}
