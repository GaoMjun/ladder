package ladder

import (
	"io"

	"github.com/golang/snappy"
)

type ConnWithSnappy struct {
	c io.ReadWriteCloser
	r *snappy.Reader
	w *snappy.Writer
}

func NewConnWithSnappy(c io.ReadWriteCloser) (conn *ConnWithSnappy) {
	conn = &ConnWithSnappy{}

	conn.c = c
	conn.r = snappy.NewReader(c)
	conn.w = snappy.NewWriter(c)
	return
}

func (self *ConnWithSnappy) Read(buf []byte) (n int, err error) {
	n, err = self.r.Read(buf)
	return
}

func (self *ConnWithSnappy) Write(buf []byte) (n int, err error) {
	n, err = self.w.Write(buf)
	return
}

func (self *ConnWithSnappy) Close() (err error) {
	err = self.w.Close()
	if err != nil {
		return
	}

	err = self.c.Close()
	return
}
