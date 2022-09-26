package ladder

import (
	"io"
	"net"
	"time"
)

type ConnWithChannel struct {
	c io.ReadWriteCloser
}

func NewConnWithChannel(c io.ReadWriteCloser) (conn *ConnWithChannel) {
	conn = &ConnWithChannel{}
	conn.c = c
	return
}

func (self *ConnWithChannel) Read(p []byte) (n int, err error) {
	n, err = self.c.Read(p)
	return
}

func (self *ConnWithChannel) Write(p []byte) (n int, err error) {
	n, err = self.c.Write(p)
	return
}

func (self *ConnWithChannel) Close() (err error) {
	err = self.c.Close()
	return
}

func (self *ConnWithChannel) LocalAddr() net.Addr {
	return nil
}

func (self *ConnWithChannel) RemoteAddr() net.Addr {
	return nil
}

func (self *ConnWithChannel) SetReadDeadline(t time.Time) error {
	return nil
}

func (self *ConnWithChannel) SetWriteDeadline(t time.Time) error {
	return nil
}

func (self *ConnWithChannel) SetDeadline(t time.Time) (err error) {
	return nil
}
