package httpstream

import (
	"io"
	"net/http/httputil"
)

type ChunckedReader struct {
	r  io.ReadCloser
	cr io.Reader
}

func NewChunckedReader(r io.ReadCloser) (cr *ChunckedReader) {
	cr = &ChunckedReader{}
	cr.r = r
	cr.cr = httputil.NewChunkedReader(r)
	return
}

func (self *ChunckedReader) Read(p []byte) (n int, err error) {
	n, err = self.cr.Read(p)
	return
}

func (self *ChunckedReader) Close() (err error) {
	err = self.r.Close()
	return
}
