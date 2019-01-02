package client

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/GaoMjun/ladder"
	"golang.org/x/crypto/ssh"
)

type TCPServer struct {
	proto    string
	addr     string
	channels *ladder.Channels
}

func NewTCPServer(proto, addr string, channels *ladder.Channels) (s *TCPServer) {
	s = &TCPServer{}
	s.proto = proto
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
		err     error
		stream  io.ReadWriteCloser
		be      *ladder.BackEnd
		sshConn ssh.Conn
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

	be, err = self.channels.GetBackEnd()
	if err != nil {
		return
	}

	sshConn = be.V.(ssh.Conn)
	log.Println(fmt.Sprint("select ", sshConn.RemoteAddr().String()))
	stream, reqs, err := sshConn.OpenChannel("", []byte(self.proto))
	if err != nil {
		return
	}
	go ssh.DiscardRequests(reqs)

	ladder.Pipe(conn, ladder.NewConnWithSnappy(stream))
}
