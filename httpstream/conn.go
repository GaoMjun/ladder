package httpstream

import (
	"io"
	"net"
	"time"
)

type Conn struct {
	r          io.ReadCloser
	w          io.WriteCloser
	remoteAddr net.Addr
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

func (self *Conn) LocalAddr() net.Addr {
	return nil
}

func (self *Conn) RemoteAddr() net.Addr {
	return self.remoteAddr
}

func (self *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (self *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (self *Conn) SetDeadline(t time.Time) error {
	return nil
}
