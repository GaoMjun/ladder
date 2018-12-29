package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/GaoMjun/ladder"
	"golang.org/x/crypto/ssh"
)

func handleConn(conn *ladder.Conn) {
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

	config, err = generateServerConfig(func(c ssh.ConnMetadata, pass []byte) (perm *ssh.Permissions, err error) {
		var (
			user   = c.User()
			passwd = string(pass)
			// session = string(c.SessionID())
		)
		if user == "fuck" && passwd == "gfw" {
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

	go handleChannels(chans)
	handleRequests(reqs)
}

func handleChannels(chans <-chan ssh.NewChannel) {
	var (
		err    error
		stream ssh.Channel
		reqs   <-chan *ssh.Request
	)
	for ch := range chans {
		stream, reqs, err = ch.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go ssh.DiscardRequests(reqs)
		go handleStream(stream)
	}
}

func handleRequests(reqs <-chan *ssh.Request) {
	for r := range reqs {
		log.Println(r)
	}
}

func handleStream(stream ssh.Channel) {
	var (
		err     error
		request *Request
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

	request, err = NewRequest(stream)
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

	if request.HttpRequest.Method == "CONNECT" {
		fmt.Fprint(stream, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		remote.Write(request.Bytes())
	}

	ladder.Pipe(stream, remote)
}
