package main

import (
	"log"
	"net"

	"github.com/GaoMjun/ladder"
	"golang.org/x/crypto/ssh"
)

func handleConn(conn net.Conn, channels *ladder.Channels) {
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
			return
		},
		ClientVersion: "SSH-2.0-ladder-1.0.0",
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
