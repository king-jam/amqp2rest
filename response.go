package mq2http

import (
	"io"
	"net/http"
)

// ResponseWriter is ...
type ResponseWriter interface {
	Status() string
	StatusCode() int
	Headers() http.Header
	Reader() io.ReadCloser
	ContentLength() int64
}

// NewResponse is ...
func NewResponse(rw ResponseWriter) (*http.Response, error) {

	resp := http.Response{
		Status:        rw.Status(),
		StatusCode:    rw.StatusCode(),
		Header:        rw.Headers(),
		Body:          rw.Reader(),
		ContentLength: rw.ContentLength(),
		Close:         false,
		Uncompressed:  true,
		Request:       nil,
	}

	return &resp, nil
}
