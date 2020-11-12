package httpstream

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/GaoMjun/goutils"
)

const (
	MAX_UP_UNIT = 1024
)

type Conn struct {
	isClient bool

	chunkedReader io.ReadCloser
	chunkedWriter io.WriteCloser

	header http.Header

	dataCh  chan []byte
	closeCh chan struct{}

	RemoteHost string

	buf []byte

	httpClient *http.Client
	serverAddr string

	writeNotWait bool
}

func (self *Conn) Read(p []byte) (n int, err error) {
	if self.isClient {
		n, err = self.chunkedReader.Read(p)
		return
	}

	if len(self.buf) > 0 {
		n = copy(p, self.buf)

		self.buf = self.buf[n:]
		return
	}

	select {
	case data := <-self.dataCh:
		if len(data) > 0 {
			n = copy(p, data)

			self.buf = data[n:]
			return
		}
	case <-self.closeCh:
		err = io.EOF
		return
	}

	return
}

func (self *Conn) Write(p []byte) (n int, err error) {
	if self.isClient {
		var (
			req  *http.Request
			resp *http.Response
		)

		if req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/%s", self.serverAddr, goutils.RandString(16)), nil); err != nil {
			return
		}
		for k, v := range self.header {
			req.Header.Set(k, v[0])
		}
		if len(p) > MAX_UP_UNIT {
			p = p[:MAX_UP_UNIT]
		}

		req.Header.Set("HTTPStream-Data", base64.StdEncoding.EncodeToString(p))

		if resp, err = self.httpClient.Do(req); err != nil {
			return
		}
		// defer resp.Body.Close()

		if self.writeNotWait {
			n = len(p)
			return
		}

		if resp.StatusCode != http.StatusNoContent {
			err = errors.New("response not ok")
			return
		}

		n = len(p)
		return
	}

	n, err = self.chunkedWriter.Write(p)
	return
}

func (self *Conn) Close() (err error) {
	if self.chunkedReader != nil {
		self.chunkedReader.Close()
	}

	if self.chunkedWriter != nil {
		self.chunkedWriter.Close()
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
