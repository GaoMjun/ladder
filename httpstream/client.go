package httpstream

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/GaoMjun/goutils"
)

func Dial(serverAddr, serverHost string, header http.Header) (conn *Conn, err error) {
	if len(serverHost) <= 0 {
		serverHost = serverAddr
	}

	if !strings.Contains(serverHost, ":") {
		serverHost += ":80"
	}

	if header == nil {
		header = http.Header{}
	}

	header.Set("Host", serverAddr)
	header.Set("HTTPStream-Key", goutils.RandString(16))

	var (
		upConn, downConn net.Conn
		reqHeader        = fmt.Sprintf("GET /%s HTTP/1.1\r\n", goutils.RandString(16))
		resp             *http.Response
	)
	defer func() {
		if err != nil {
			if upConn != nil {
				upConn.Close()
			}

			if downConn != nil {
				downConn.Close()
			}
		}
	}()

	if upConn, err = net.Dial("tcp", serverHost); err != nil {
		return
	}

	if downConn, err = net.Dial("tcp", serverHost); err != nil {
		return
	}

	for k, v := range header {
		reqHeader += fmt.Sprintf("%s: %s\r\n", k, v[0])
	}
	reqHeader += "\r\n"

	if _, err = downConn.Write([]byte(reqHeader)); err != nil {
		return
	}

	if resp, err = http.ReadResponse(bufio.NewReader(downConn), nil); err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New("response not ok")
		return
	}

	conn = &Conn{}
	conn.isClient = true
	conn.upConn = upConn
	conn.downConn = downConn
	conn.chunkedReader = resp.Body
	conn.header = header

	return
}
