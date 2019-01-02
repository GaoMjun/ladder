package ladder

import (
	"net"
	"time"
)

type ConnWithXor struct {
	k []byte
	c net.Conn
}

func NewConnWithXor(c net.Conn, key []byte) (conn *ConnWithXor) {
	conn = &ConnWithXor{}
	conn.c = c
	conn.k = key
	return
}

func (self *ConnWithXor) Read(p []byte) (n int, err error) {
	n, err = self.c.Read(p)
	if err != nil {
		return
	}

	o := xor(p[:n], self.k)
	copy(p, o)
	return
}

func (self *ConnWithXor) Write(p []byte) (n int, err error) {
	o := xor(p, self.k)
	return self.c.Write(o)
}

func (self *ConnWithXor) Close() (err error) {
	err = self.c.Close()
	return
}

func (self *ConnWithXor) LocalAddr() net.Addr {
	return self.c.LocalAddr()
}

func (self *ConnWithXor) RemoteAddr() net.Addr {
	return self.c.RemoteAddr()
}

func (self *ConnWithXor) SetReadDeadline(t time.Time) error {
	return self.c.SetReadDeadline(t)
}

func (self *ConnWithXor) SetWriteDeadline(t time.Time) error {
	return self.c.SetWriteDeadline(t)
}

func (self *ConnWithXor) SetDeadline(t time.Time) (err error) {
	err = self.c.SetDeadline(t)
	return
}
