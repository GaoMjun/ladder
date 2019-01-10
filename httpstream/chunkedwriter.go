package httpstream

import (
	"io"
	"net/http/httputil"
)

type ChunckedWriter struct {
	w  io.WriteCloser
	cw io.Writer
}

func NewChunckedWriter(w io.WriteCloser) (cw *ChunckedWriter) {
	cw = &ChunckedWriter{}
	cw.w = w
	cw.cw = httputil.NewChunkedWriter(w)
	return
}

func (self *ChunckedWriter) Write(p []byte) (n int, err error) {
	n, err = self.cw.Write(p)
	return
}

func (self *ChunckedWriter) Close() (err error) {
	err = self.w.Close()
	return
}
