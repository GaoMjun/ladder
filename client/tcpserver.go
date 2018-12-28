package main

import (
	"io"
	"log"
	"net"

	"github.com/GaoMjun/ladder"
)

type TCPServer struct {
	addr     string
	channels *ladder.Channels
}

func NewTCPServer(addr string, channels *ladder.Channels) (s *TCPServer) {
	s = &TCPServer{}
	s.addr = addr
	s.channels = channels
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
		err    error
		stream io.ReadWriteCloser
	)
	defer func() {
		conn.Close()
		if stream != nil {
			stream.Close()
		}
		if err != nil {
			log.Println(err)
		}
	}()

	stream, err = self.channels.CreateStream()
	if err != nil {
		return
	}

	ladder.Pipe(conn, stream)
}
