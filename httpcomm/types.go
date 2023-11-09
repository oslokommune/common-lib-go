package httpcomm

import (
	"fmt"
	"io"
)

type HTTPRequest struct {
	Body    io.Reader
	Token   *string
	Headers map[string]string
	Url     string
	Method  string
	Tracing bool
}

type HTTPResponse struct {
	Body       string
	StatusCode int
}

type HTTPError struct {
	Body       string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP-status: %d, Message: %s", e.StatusCode, e.Body)
}
