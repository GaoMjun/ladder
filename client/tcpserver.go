package client

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/GaoMjun/ladder"
)

type TCPServer struct {
	proto         string
	addr          string
	channels      *ladder.Channels
	streamManager *ladder.StreamManager
}

func NewTCPServer(proto, addr string, channels *ladder.Channels, streamManager *ladder.StreamManager) (s *TCPServer) {
	s = &TCPServer{}
	s.proto = proto
	s.addr = addr
	s.channels = channels
	s.streamManager = streamManager
	return
}

func (self *TCPServer) Run() (err error) {
	var (
		l net.Listener
	)

	l, err = net.Listen("tcp", self.addr)
	if err != nil {
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go self.handleConn(conn)
	}

	return
}

func (self *TCPServer) handleConn(conn net.Conn) {
	var (
		err error

		be           *ladder.BackEnd
		request      *ladder.Request
		streamID     uint16
		wc           io.WriteCloser
		originHeader string
	)
	defer func() {
		conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	streamID, err = self.streamManager.GenerateStreamID()
	if err != nil {
		return
	}

	be, err = self.channels.GetBackEnd()
	if err != nil {
		return
	}

	channel := be.V.(*Channel)
	token, _ := ladder.GenerateToken(channel.user, channel.pass)

	if self.proto == "http" {
		request, err = ladder.NewRequest(conn)
		if err != nil {
			return
		}

		if request.HttpRequest.Method == "CONNECT" {
			fmt.Fprint(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
		}

		originHeader, err = ladder.EncryptHeader(request.Dump(), channel.user, channel.pass)
		if err != nil {
			return
		}

		fakeRequest := ladder.NewFakeRequest(channel.host, token, originHeader, fmt.Sprint(streamID))
		c, err := fakeRequest.Do()
		if err != nil {
			return
		}
		wc = c.(io.WriteCloser)
	}

	rc := self.streamManager.NewReceiveStream(streamID)

	stream := ladder.NewConnCombine(wc, rc)

	ladder.Pipe(conn, stream)
}
