package ladder

import (
	"fmt"
	"net"
	"time"
)

type FakeRequest struct {
	fakeHeader string
	dialer     *net.Dialer
	host       string
}

func NewFakeRequest(host, token, originHeader, streamID string) (request *FakeRequest) {
	request = &FakeRequest{}
	request.host = host

	request.fakeHeader = "GET / HTTP/1.1\r\nHost: www.baidu.com\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3614.0 Safari/537.36\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\nCache-Control: no-cache\r\n"
	if len(token) > 0 {
		request.fakeHeader += fmt.Sprint("Token: ", token, "\r\n")
	}
	if len(originHeader) > 0 {
		request.fakeHeader += fmt.Sprint("Header: ", originHeader, "\r\n")
	}
	if len(streamID) > 0 {
		request.fakeHeader += fmt.Sprint("Sid: ", streamID, "\r\n")
	}
	request.fakeHeader += "\r\n"

	request.dialer = &net.Dialer{Timeout: time.Second * 13}
	return
}

func (self *FakeRequest) Do() (conn net.Conn, err error) {
	conn, err = self.dialer.Dial("tcp", self.host)
	if err != nil {
		return
	}

	_, err = fmt.Fprint(conn, self.fakeHeader)
	if err != nil {
		return
	}
	return
}
