package http

import (
	"bytes"
	"net/http"
)

var _ http.ResponseWriter = (*ResponseWriter)(nil)

type ResponseWriter struct {
	headers    http.Header
	statusCode int
	buffer     *bytes.Buffer
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		headers: make(http.Header),
		buffer:  bytes.NewBuffer(nil),
	}
}

func (r *ResponseWriter) Header() http.Header {
	return r.headers
}

func (r *ResponseWriter) Write(bytes []byte) (int, error) {
	return r.buffer.Write(bytes)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}
