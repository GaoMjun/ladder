package client

import (
	"fmt"
	"io"
	"log"
	"net"

	"ladder"

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

	if l, err = net.Listen("tcp", self.addr); err != nil {
		return
	}
	log.Println("local", self.proto, "server listen at", self.addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go self.handleConn(conn)
	}
}

func (self *TCPServer) handleConn(conn net.Conn) {
	var (
		err          error
		stream       io.ReadWriteCloser
		be           *ladder.BackEnd
		sshConn      ssh.Conn
		request      *ladder.Request
		snappyStream io.ReadWriteCloser
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

	channel := be.V.(*Channel)
	sshConn = channel.conn
	log.Println(fmt.Sprint("select ", sshConn.RemoteAddr().String()))
	extraData := []byte(self.proto)
	if channel.comp == true {
		extraData = append(extraData, []byte(":comp")...)
	}
	stream, reqs, err := sshConn.OpenChannel("", extraData)
	if err != nil {
		return
	}
	if channel.comp == true {
		snappyStream = ladder.NewConnWithSnappy(stream)
	} else {
		snappyStream = stream
	}

	go ssh.DiscardRequests(reqs)

	if self.proto == "http" {
		request, err = ladder.NewRequest(conn)
		if err != nil {
			return
		}

		if request.HttpRequest.Method == "CONNECT" {
			fmt.Fprint(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
		}

		snappyStream.Write(request.Bytes())
	}

	ladder.Pipe(conn, snappyStream)
}
