package ladder

import (
	"errors"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	wsConn *websocket.Conn
	buffer []byte
}

func NewConn(wsConn *websocket.Conn) (conn *Conn) {
	conn = &Conn{}
	conn.wsConn = wsConn
	return
}

func (self *Conn) Read(buf []byte) (n int, err error) {
	var (
		t    int
		msg  []byte
		size int
	)

	if len(self.buffer) > 0 {
		if len(self.buffer) <= len(buf) {
			buf = self.buffer
			n = len(self.buffer)
		} else {
			buf = self.buffer[:len(buf)]
			n = len(buf)
			self.buffer = self.buffer[len(buf):]
		}
	}

	t, msg, err = self.wsConn.ReadMessage()
	if err != nil {
		return
	}
	if t != websocket.BinaryMessage {
		err = errors.New("websocket non-binary msg")
		return
	}

	size = len(msg)
	if size <= len(buf) {
		buf = msg[:size]
		n = size
	} else {
		buf = msg[:len(buf)]
		n = len(buf)
		self.buffer = msg[len(buf):]
	}
	return
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
