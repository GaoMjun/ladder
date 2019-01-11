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
)

type Dialer struct {
	NetDial   func(network, addr string) (net.Conn, error)
	UpNetDial func(network, addr string) (net.Conn, error)
	Timeout   time.Duration
	UpHost    string

	upConn   net.Conn
	downConn net.Conn
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

	if header == nil {
		header = http.Header{}
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
		wg         = &sync.WaitGroup{}
		w          io.WriteCloser
		r          io.ReadCloser
		remoteAddr net.Addr
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

		w, _, err = self.openUp(u, header)
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

		r, remoteAddr, err = self.openDown(u, header)
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
	conn.remoteAddr = remoteAddr
	return
}

func (self *Dialer) openUp(u *url.URL, header http.Header) (w io.WriteCloser, remoteAddr net.Addr, err error) {
	var (
		request       *http.Request
		requestHeader string
		bs            []byte
		netConn       net.Conn
		host          = u.Host
	)

	if len(self.UpHost) > 0 {
		host = self.UpHost
	}

	requestHeader += fmt.Sprintf("POST %s HTTP/1.1\r\n", u.Path)
	requestHeader += fmt.Sprintf("Host: %s\r\n", host)
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

	if self.UpNetDial != nil {
		netConn, err = self.UpNetDial("tcp", host)
	} else {
		netConn, err = self.NetDial("tcp", host)
	}
	if err != nil {
		return
	}

	remoteAddr = netConn.RemoteAddr()

	self.upConn = netConn

	_, err = self.upConn.Write(bs)
	if err != nil {
		return
	}

	w = NewChunckedWriter(self.upConn)
	return
}

func (self *Dialer) openDown(u *url.URL, header http.Header) (r io.ReadCloser, remoteAddr net.Addr, err error) {
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
	remoteAddr = netConn.RemoteAddr()

	self.downConn = netConn

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
