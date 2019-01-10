package ladder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

type HTTPStream struct {
	r                io.ReadCloser
	w                io.WriteCloser
	UpCustomHeader   string
	DownCustomHeader string
	Timeout          time.Duration
	Host             string
	dialer           *net.Dialer
	rawurl           string
	upConn           *ConnWithTimeout
	downConn         *ConnWithTimeout
	headers          []string
	path             string
}

func OpenStream(u string) (s *HTTPStream, err error) {
	s = &HTTPStream{}
	err = s.Open(u)
	return
}

func (self *HTTPStream) AddHeader(k, v string) {
	self.headers = append(self.headers, fmt.Sprint(k, ": ", v, "\r\n"))
}

func (self *HTTPStream) Open(rawurl string) (err error) {
	if self.Timeout == 0 {
		self.Timeout = time.Second * 30
	}

	self.dialer = &net.Dialer{Timeout: self.Timeout}

	u, err := url.Parse(rawurl)
	if err != nil {
		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New(fmt.Sprint("not support protocol ", u.Scheme))
		return
	}

	if len(self.Host) <= 0 {
		port := u.Port()

		switch u.Scheme {
		case "http":
			if len(port) <= 0 {
				port = "80"
			}
		case "https":
			if len(port) <= 0 {
				port = "443"
			}
		}

		self.Host = fmt.Sprint(u.Hostname(), ":", port)
	}

	self.path = u.Path
	if len(self.path) <= 0 {
		self.path = "/"
	}

	var (
		wg = &sync.WaitGroup{}
		w  io.WriteCloser
		r  io.ReadCloser
	)
	defer func() {
		if err != nil {
			if w != nil {
				w.Close()
			}
			if r != nil {
				r.Close()
			}
		}
	}()

	wg.Add(1)
	go func() {
		var (
			err error
		)
		defer func() {
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}()

		w, err = self.openUp()
		if err != nil {
			return
		}
	}()

	wg.Add(1)
	go func() {
		var (
			err error
		)
		defer func() {
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}()

		r, err = self.openDown()
		if err != nil {
			return
		}
	}()

	wg.Wait()

	if w == nil || r == nil {
		err = errors.New("open stream failed")
		return
	}

	self.w = w
	self.r = r

	return
}

func (self *HTTPStream) Read(p []byte) (n int, err error) {
	n, err = self.r.Read(p)
	return
}

func (self *HTTPStream) Write(p []byte) (n int, err error) {
	n, err = self.w.Write(p)
	return
}

func (self *HTTPStream) Close() (err error) {
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

func (self *HTTPStream) openUp() (w io.WriteCloser, err error) {
	var (
		request *http.Request
		bs      []byte
	)

	if len(self.UpCustomHeader) <= 0 {
		self.UpCustomHeader += fmt.Sprintf("POST %s HTTP/1.1\r\n", self.path)
		self.UpCustomHeader += fmt.Sprintf("Host: %s\r\n", self.Host)
		self.UpCustomHeader += fmt.Sprintf("Content-Type: %s\r\n", "octet-stream")
		self.UpCustomHeader += fmt.Sprintf("Transfer-Encoding: %s\r\n", "chunked")

		for _, h := range self.headers {
			self.UpCustomHeader += h
		}

		self.UpCustomHeader += "\r\n"
	}

	request, err = http.ReadRequest(bufio.NewReader(strings.NewReader(self.UpCustomHeader)))
	if err != nil {
		return
	}

	bs, err = httputil.DumpRequest(request, false)
	if err != nil {
		return
	}
	// fmt.Println(string(bs))

	conn, err := self.dialer.Dial("tcp", self.Host)
	if err != nil {
		return
	}

	self.upConn = &ConnWithTimeout{Conn: conn, Timeout: self.Timeout}

	_, err = self.upConn.Write(bs)
	if err != nil {
		return
	}

	w = NewChunckedWriter(self.upConn)
	return
}

func (self *HTTPStream) openDown() (r io.ReadCloser, err error) {
	var (
		request  *http.Request
		response *http.Response
		bs       []byte
	)

	if len(self.DownCustomHeader) <= 0 {
		self.DownCustomHeader += fmt.Sprintf("GET %s HTTP/1.1\r\n", self.path)
		self.DownCustomHeader += fmt.Sprintf("Host: %s\r\n", self.Host)

		for _, h := range self.headers {
			self.DownCustomHeader += h
		}

		self.DownCustomHeader += "\r\n"
	}

	request, err = http.ReadRequest(bufio.NewReader(strings.NewReader(self.DownCustomHeader)))
	if err != nil {
		return
	}

	bs, err = httputil.DumpRequest(request, false)
	if err != nil {
		return
	}
	// fmt.Println(string(bs))

	conn, err := self.dialer.Dial("tcp", self.Host)
	if err != nil {
		return
	}

	self.downConn = &ConnWithTimeout{Conn: conn, Timeout: self.Timeout}

	_, err = self.downConn.Write(bs)
	if err != nil {
		return
	}

	response, err = http.ReadResponse(bufio.NewReader(self.downConn), nil)
	if err != nil {
		return
	}

	// bs, err = httputil.DumpResponse(response, false)
	// fmt.Println(string(bs))

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	r = NewChunckedReader(self.downConn)
	return
}
