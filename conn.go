package ladder

import (
	"io"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	wsConn *websocket.Conn
	buffer []byte
	rbuf   []byte
}

func NewConn(wsConn *websocket.Conn) (conn *Conn) {
	conn = &Conn{}
	conn.wsConn = wsConn
	conn.rbuf = make([]byte, ^uint16(0))
	return
}

func (self *Conn) Read(buf []byte) (n int, err error) {
TRY:
	if len(self.buffer) > 0 {
		if len(self.buffer) >= len(buf) {
			n = copy(buf, self.buffer)
			self.buffer = self.buffer[n:]
			return
		}

		n = copy(buf, self.buffer)
		self.buffer = nil
		return
	}

	var (
		size = 0
	)
	if size, err = readMessage(self.wsConn, self.rbuf); err != nil {
		return
	}

	self.buffer = append(self.buffer, self.rbuf[:size]...)
	goto TRY
}

func (self *Conn) Write(buf []byte) (n int, err error) {
	err = self.wsConn.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		return
	}
	n = len(buf)
	return
}

func (self *Conn) Close() error {
	return self.wsConn.Close()
}

func (self *Conn) LocalAddr() net.Addr {
	return self.wsConn.LocalAddr()
}

func (self *Conn) RemoteAddr() net.Addr {
	return self.wsConn.RemoteAddr()
}

func (self *Conn) SetReadDeadline(t time.Time) error {
	return self.wsConn.SetReadDeadline(t)
}

func (self *Conn) SetWriteDeadline(t time.Time) error {
	return self.wsConn.SetWriteDeadline(t)
}

func (self *Conn) SetDeadline(t time.Time) (err error) {
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

func readMessage(wsConn *websocket.Conn, buf []byte) (n int, err error) {
TRY:
	var (
		r           io.Reader
		messageType int
	)
	if messageType, r, err = wsConn.NextReader(); err != nil {
		return
	}
	if messageType != websocket.BinaryMessage {
		goto TRY
	}

	if n, err = io.ReadFull(r, buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = nil
			return
		}
	}

	return
}
