package httpstream

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

	"github.com/GaoMjun/ladder"
)

type Dialer struct {
	NetDial func(network, addr string) (net.Conn, error)
	Timeout time.Duration

	upConn   *ladder.ConnWithTimeout
	downConn *ladder.ConnWithTimeout
}

func (self *Dialer) Dial(rawurl string, header http.Header) (conn *Conn, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New(fmt.Sprint("not support protocol ", u.Scheme))
		return
	}

	if len(u.Port()) <= 0 {
		switch u.Scheme {
		case "http":
			u.Host = fmt.Sprint(u.Hostname(), ":80")
		case "https":
			u.Host = fmt.Sprint(u.Hostname(), ":443")
		}
	}

	if v, ok := header["Host"]; ok {
		u.Host = v[0]
	}

	if len(u.Path) <= 0 {
		u.Path = "/"
	}

	if self.NetDial == nil {
		self.NetDial = func(network, addr string) (net.Conn, error) {
			d := &net.Dialer{Timeout: self.Timeout}
			return d.Dial(network, addr)
		}
	}

	if self.Timeout <= 0 {
		self.Timeout = time.Second * 30
	}

	key, err := generateChallengeKey()
	if err != nil {
		return
	}
	header["HTTPStream-Key"] = []string{key}

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

		w, err = self.openUp(u, header)
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

		r, err = self.openDown(u, header)
		if err != nil {
			return
		}
	}()

	wg.Wait()

	if w == nil || r == nil {
		err = errors.New("open stream failed")
		return
	}

	conn = &Conn{}
	conn.w = w
	conn.r = r
	return
}

func (self *Dialer) openUp(u *url.URL, header http.Header) (w io.WriteCloser, err error) {
	var (
		request       *http.Request
		requestHeader string
		bs            []byte
		netConn       net.Conn
	)

	requestHeader += fmt.Sprintf("POST %s HTTP/1.1\r\n", u.Path)
	requestHeader += fmt.Sprintf("Host: %s\r\n", u.Host)
	requestHeader += fmt.Sprintf("Content-Type: %s\r\n", "octet-stream")
	requestHeader += fmt.Sprintf("Transfer-Encoding: %s\r\n", "chunked")

	for k, vs := range header {
		if k == "Host" || k == "Content-Type" || k == "Transfer-Encoding" {
			continue
		}

		requestHeader += fmt.Sprintf("%s: %s\r\n", k, vs[0])
	}
	requestHeader += "\r\n"

	request, err = http.ReadRequest(bufio.NewReader(strings.NewReader(requestHeader)))
	if err != nil {
		return
	}

	bs, err = httputil.DumpRequest(request, false)
	if err != nil {
		return
	}
	// fmt.Println(string(bs))

	netConn, err = self.NetDial("tcp", u.Host)
	if err != nil {
		return
	}

	self.upConn = &ladder.ConnWithTimeout{Conn: netConn, Timeout: self.Timeout}

	_, err = self.upConn.Write(bs)
	if err != nil {
		return
	}

	w = NewChunckedWriter(self.upConn)
	return
}

func (self *Dialer) openDown(u *url.URL, header http.Header) (r io.ReadCloser, err error) {
	var (
		request       *http.Request
		requestHeader string
		response      *http.Response
		bs            []byte
		netConn       net.Conn
	)

	requestHeader += fmt.Sprintf("GET %s HTTP/1.1\r\n", u.Path)
	requestHeader += fmt.Sprintf("Host: %s\r\n", u.Host)

	for k, vs := range header {
		if k == "Host" {
			continue
		}

		requestHeader += fmt.Sprintf("%s: %s\r\n", k, vs[0])
	}
	requestHeader += "\r\n"

	request, err = http.ReadRequest(bufio.NewReader(strings.NewReader(requestHeader)))
	if err != nil {
		return
	}

	bs, err = httputil.DumpRequest(request, false)
	if err != nil {
		return
	}
	// fmt.Println(string(bs))

	netConn, err = self.NetDial("tcp", u.Host)
	if err != nil {
		return
	}

	self.downConn = &ladder.ConnWithTimeout{Conn: netConn, Timeout: self.Timeout}

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
