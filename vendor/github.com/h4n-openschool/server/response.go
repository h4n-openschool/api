package server

import (
	"bytes"
	"io"
	"net/http"
)

func NewResponse() http.Response {
	return http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
	}
}

func SetBody(r http.Response, b []byte) http.Response {
	r.Body = io.NopCloser(bytes.NewReader(b))
	return r
}
