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

func (self nopHttpResponseWriteCloser) Close() (err error) {
	// if self.w == nil {
	// 	return
	// }

	// v, ok := self.w.(http.Hijacker)
	// if !ok {
	// 	return
	// }

	// var conn net.Conn
	// if conn, _, err = v.Hijack(); err != nil {
	// 	return
	// }
	// err = conn.Close()
	return
}

func NopHttpResponseWriteCloser(w http.ResponseWriter) (wc io.WriteCloser) {
	return &nopHttpResponseWriteCloser{w}
}
