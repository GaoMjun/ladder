package httpstream

import (
	"io"
	"net/http"
)

type nopHttpResponseWriteCloser struct {
	w http.ResponseWriter
}

func (self nopHttpResponseWriteCloser) Write(p []byte) (n int, err error) {
	if n, err = self.w.Write(p); err != nil {
		return
	}

	self.w.(http.Flusher).Flush()
	return
}

func (self nopHttpResponseWriteCloser) Close() error {
	return nil
}

func NopHttpResponseWriteCloser(w http.ResponseWriter) (wc io.WriteCloser) {
	return &nopHttpResponseWriteCloser{w}
}
