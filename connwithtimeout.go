package ladder

import (
	"net"
	"time"
)

type ConnWithTimeout struct {
	Timeout time.Duration
	Conn    net.Conn
}

func NewConnWithTimeout(conn net.Conn) (c *ConnWithTimeout) {
	c = &ConnWithTimeout{}
	c.Conn = conn
	c.Timeout = time.Second * 30
	return
}

func (self *ConnWithTimeout) Read(p []byte) (n int, err error) {
	if self.Timeout > 0 {
		self.Conn.SetReadDeadline(time.Now().Add(self.Timeout))
	} else {
		self.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	}
	return self.Conn.Read(p)
}

func (self *ConnWithTimeout) Write(p []byte) (n int, err error) {
	if self.Timeout > 0 {
		self.Conn.SetWriteDeadline(time.Now().Add(self.Timeout))
	} else {
		self.Conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	}
	return self.Conn.Write(p)
}

func (self *ConnWithTimeout) Close() (err error) {
	err = self.Conn.Close()
	return
}
