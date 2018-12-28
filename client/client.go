package main

import (
	"log"
	"net"

	"github.com/GaoMjun/ladder"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var (
		err  error
		conn *websocket.Conn

		tcpServer *TCPServer
		channels  = &ladder.Channels{}
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	conn, _, err = websocket.DefaultDialer.Dial("ws://127.0.0.1:8888/", nil)
	if err != nil {
		return
	}

	handleConn(ladder.NewConn(conn), channels)

	tcpServer = NewTCPServer("127.0.0.1:9999", channels)
	err = tcpServer.Run()
}

func handleConn(conn *ladder.Conn, channels *ladder.Channels) {
	var (
		err     error
		config  *ssh.ClientConfig
		sshConn ssh.Conn
		chans   <-chan ssh.NewChannel
		reqs    <-chan *ssh.Request

		backend *ladder.BackEnd
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	config = &ssh.ClientConfig{
		User: "fuck",
		Auth: []ssh.AuthMethod{ssh.Password("gfw")},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) (err error) {
			log.Println(hostname, remote)
			return
		},
	}

	sshConn, chans, reqs, err = ssh.NewClientConn(conn, "", config)
	if err != nil {
		return
	}

	go ssh.DiscardRequests(reqs)
	go func() {
		for ch := range chans {
			ch.Reject(ssh.Prohibited, "Tunnels disallowed")
		}
	}()

	backend = ladder.NewBackEnd(sshConn)

	channels.AddBackEnd(backend)

	err = sshConn.Wait()

	channels.DelBackEnd(backend)
}
