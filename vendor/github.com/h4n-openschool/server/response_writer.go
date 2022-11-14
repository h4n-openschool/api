package server

import "net/http"

type OSResponseWriter struct {
	body       []byte
	statusSet  bool
	statusCode int
	header     http.Header
}

func NewOSResponseWriter() *OSResponseWriter {
	return &OSResponseWriter{
		statusSet: false,
		header:    http.Header{},
	}
}

func (w *OSResponseWriter) Header() http.Header {
	return w.header
}

func (w *OSResponseWriter) Write(b []byte) (int, error) {
	if !w.statusSet {
		w.WriteHeader(200)
	}

	w.body = b

	return 0, nil
}

func (w *OSResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.statusSet = true
}
