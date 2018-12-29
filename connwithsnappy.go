package ladder

import (
	"net"
	"time"

	"github.com/golang/snappy"
)

type ConnWithSnappy struct {
	c net.Conn
	r *snappy.Reader
	w *snappy.Writer
}

func NewConnWithSnappy(c net.Conn) (conn *ConnWithSnappy) {
	conn = &ConnWithSnappy{}

	conn.c = c
	conn.r = snappy.NewReader(c)
	conn.w = snappy.NewBufferedWriter(c)
	return
}

func (self *ConnWithSnappy) Read(buf []byte) (n int, err error) {
	n, err = self.r.Read(buf)
	return
}

func (self *ConnWithSnappy) Write(buf []byte) (n int, err error) {
	n, err = self.w.Write(buf)
	if err != nil {
		return
	}
	err = self.w.Flush()
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

func (self *ConnWithSnappy) LocalAddr() net.Addr {
	return self.c.LocalAddr()
}

func (self *ConnWithSnappy) RemoteAddr() net.Addr {
	return self.c.RemoteAddr()
}

func (self *ConnWithSnappy) SetReadDeadline(t time.Time) error {
	return self.c.SetReadDeadline(t)
}

func (self *ConnWithSnappy) SetWriteDeadline(t time.Time) error {
	return self.c.SetWriteDeadline(t)
}

func (self *ConnWithSnappy) SetDeadline(t time.Time) (err error) {
	err = self.SetReadDeadline(t)
	if err != nil {
		return
	}
	err = self.SetWriteDeadline(t)
	if err != nil {
		return
	}
	return
}
