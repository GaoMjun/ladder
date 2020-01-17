package httpstream

import (
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
	header.Set("HTTPStream-Key", goutils.RandString(16))

	var (
		req  *http.Request
		resp *http.Response

		httpClient = &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        0,
				MaxIdleConnsPerHost: 128,
				MaxConnsPerHost:     0,
				Dial: func(network, addr string) (conn net.Conn, err error) {
					if serverHost != "" {
						addr = serverHost
					}

					return net.Dial(network, addr)
				},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	)

	if req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/%s", serverAddr, goutils.RandString(16)), nil); err != nil {
		return
	}
	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	if resp, err = httpClient.Do(req); err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New("response not ok")
		return
	}

	conn = &Conn{}
	conn.isClient = true
	conn.chunkedReader = resp.Body
	conn.header = header
	conn.httpClient = httpClient
	conn.serverAddr = serverAddr

	return
}
