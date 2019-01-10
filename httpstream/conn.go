package httpstream

import (
	"io"
)

type Conn struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (self *Conn) Read(p []byte) (n int, err error) {
	n, err = self.r.Read(p)
	return
}

func (self *Conn) Write(p []byte) (n int, err error) {
	n, err = self.w.Write(p)
	return
}

func (self *Conn) Close() (err error) {
	err = self.r.Close()
	if err != nil {
		return
	}

	err = self.w.Close()
	if err != nil {
		return
	}
	return
}
