package httpcomm

import (
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
	Error      bool
}

type HttpError struct {
	Message    string
	StatusCode int
}
