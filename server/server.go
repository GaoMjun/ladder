package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/GaoMjun/tunsocks/iptransparent"

	"github.com/GaoMjun/ladder"
	"golang.org/x/crypto/ssh"
)

func handleConn(conn net.Conn, user, pass string, comp bool) {
	var (
		err     error
		config  *ssh.ServerConfig
		sshConn *ssh.ServerConn
		chans   <-chan ssh.NewChannel
		reqs    <-chan *ssh.Request
	)
	defer func() {
		conn.Close()
		if sshConn != nil {
			sshConn.Close()
		}
		if err != nil {
			log.Println(err)
		}
	}()

	config, err = generateServerConfig(func(c ssh.ConnMetadata, passwd []byte) (perm *ssh.Permissions, err error) {
		var (
			u = c.User()
			p = string(passwd)
			// session = string(c.SessionID())
		)
		if u == user && p == pass {
			return
		}
		err = errors.New("Invalid auth")
		return
	})

	if err != nil {
		return
	}

	sshConn, chans, reqs, err = ssh.NewServerConn(conn, config)
	if err != nil {
		return
	}

	go handleChannels(chans, comp)
	handleRequests(reqs)
}

func handleChannels(chans <-chan ssh.NewChannel, comp bool) {
	var (
		err          error
		stream       ssh.Channel
		reqs         <-chan *ssh.Request
		snappyStream io.ReadWriteCloser
	)

	for ch := range chans {
		stream, reqs, err = ch.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		if comp == true {
			snappyStream = ladder.NewConnWithSnappy(stream)
		} else {
			snappyStream = stream
		}

		go ssh.DiscardRequests(reqs)

		proto := string(ch.ExtraData())
		switch proto {
		case "http":
			go handleHTTP(snappyStream)
		case "socks":
			go handleSocks(ladder.NewConnWithChannel(snappyStream))
		case "iptransparent":
			go handleIPTransparent(snappyStream)
		default:
			log.Println(fmt.Sprint("not support protocol ", proto))
		}
	}
}

var socksServer = ladder.NewSocks5Server()
var iptransparentServer = &iptransparent.Server{}

func handleSocks(conn net.Conn) {
	err := socksServer.ServeConn(conn)
	if err != nil {
		log.Println(err)
	}
}

func handleIPTransparent(conn io.ReadWriteCloser) {
	err := iptransparentServer.ServeConn(conn)
	if err != nil {
		log.Println(err)
	}
}

func handleRequests(reqs <-chan *ssh.Request) {
	for r := range reqs {
		switch r.Type {
		case "ping":
			err := r.Reply(true, nil)
			if err != nil {
				log.Println(err)
			}
			log.Println("ping keepalive")
		default:
			log.Println(r)
		}
	}
}

func handleHTTP(stream io.ReadWriteCloser) {
	var (
		err     error
		request *ladder.Request
		address string
		conn    net.Conn
		remote  *ladder.ConnWithTimeout
		dialer  = &net.Dialer{Timeout: time.Second * 3}
	)
	defer func() {
		stream.Close()
		if remote != nil {
			remote.Close()
		}
		if err != nil {
			log.Println(err)
		}
	}()

	request, err = ladder.NewRequest(stream)
	if err != nil {
		return
	}

	address = request.HttpRequest.Host
	if strings.Index(address, ":") == -1 {
		address = address + ":80"
	}
	log.Println(address)

	conn, err = dialer.Dial("tcp", address)
	if err != nil {
		return
	}
	remote = ladder.NewConnWithTimeout(conn)

	if request.HttpRequest.Method != "CONNECT" {
		remote.Write(request.Bytes())
	}

	ladder.Pipe(stream, remote)
}
