package mq2http

import (
	"fmt"
	"io"
	"net/http"
)

// RequestWriter is ...
type RequestWriter interface {
	Headers() map[string]string
	Method() string
	URL() string
	Reader() io.Reader
}

// NewRequest is ...
func NewRequest(rw RequestWriter) (*http.Request, error) {
	req, err := http.NewRequest(rw.Method(), rw.URL(), rw.Reader())
	if err != nil {
		return req, fmt.Errorf("Failed to create base request")
	}
	for k, v := range rw.Headers() {
		req.Header.Set(k, v)
	}
	return req, nil
}
