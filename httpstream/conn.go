package httpstream

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/GaoMjun/goutils"
)

type Conn struct {
	isClient         bool
	upConn, downConn net.Conn

	chunkedReader io.ReadCloser
	chunkedWriter io.WriteCloser

	header http.Header

	dataCh chan []byte
}

func (self *Conn) Read(p []byte) (n int, err error) {
	if self.isClient {
		n, err = self.chunkedReader.Read(p)
		return
	}

	data, ok := <-self.dataCh
	if !ok {
		err = io.EOF
		return
	}

	if len(data) > 0 {
		if len(data) > len(p) {
			err = errors.New("read buffer not enough")
			return
		}

		n = copy(p, data)
	}
	return
}

func (self *Conn) Write(p []byte) (n int, err error) {
	if self.isClient {
		n = len(p)

		header := fmt.Sprintf("POST /%s HTTP/1.1\r\n", goutils.RandString(16))
		header += fmt.Sprintf("Connection: keep-alive\r\n")
		header += fmt.Sprintf("Cache-Control: no-cache\r\n")
		header += fmt.Sprintf("Content-Type: application/octet-stream\r\n")
		header += fmt.Sprintf("Content-Length: %d\r\n", n)
		for k, v := range self.header {
			header += fmt.Sprintf("%s: %s\r\n", k, v[0])
		}
		header += "\r\n"

		if _, err = self.upConn.Write([]byte(header)); err != nil {
			return
		}

		if _, err = self.upConn.Write(p); err != nil {
			return
		}

		var resp *http.Response
		if resp, err = http.ReadResponse(bufio.NewReader(self.upConn), nil); err != nil {
			return
		}

		if resp.StatusCode != http.StatusNoContent {
			err = errors.New("response not ok")
			return
		}
		return
	}

	n, err = self.chunkedWriter.Write(p)
	return
}

func (self *Conn) Close() (err error) {
	if self.upConn != nil {
		self.upConn.Close()
	}

	if self.downConn != nil {
		self.downConn.Close()
	}

	return
}

func (self *Conn) LocalAddr() net.Addr {
	return nil
}

func (self *Conn) RemoteAddr() net.Addr {
	return nil
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
