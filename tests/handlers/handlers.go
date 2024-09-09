package handlers_test

import (
	"errors"
	"net/http"
)

type DummyResponseWriter struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

func NewDummyResponseWriter() *DummyResponseWriter {
	return &DummyResponseWriter{
		StatusCode: 200,
		Body:       nil,
		Headers:    http.Header{},
	}
}

func (d *DummyResponseWriter) Header() http.Header {
	return d.Headers
}

func (d *DummyResponseWriter) Write(b []byte) (int, error) {
	if d.Body != nil {
		return 0, errors.New("body has already been written")
	}

	d.Body = b
	return len(b), nil
}

func (d *DummyResponseWriter) WriteHeader(statusCode int) {
	d.StatusCode = statusCode
}
