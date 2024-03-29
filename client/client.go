package client

import (
	"log"
	"net"
	"time"

	"github.com/GaoMjun/goutils"

	"github.com/GaoMjun/ladder"
	"golang.org/x/crypto/ssh"
)

func handleConn(user, pass string, comp bool, conn net.Conn, channels *ladder.Channels, successFunc func()) {
	var (
		err     error
		config  *ssh.ClientConfig
		sshConn ssh.Conn
		chans   <-chan ssh.NewChannel
		reqs    <-chan *ssh.Request

		backend  *ladder.BackEnd
		repeater *goutils.Repeater
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	config = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) (err error) {
			return
		},
		ClientVersion: "SSH-2.0-ladder-1.0.0",
	}

	sshConn, chans, reqs, err = ssh.NewClientConn(conn, "", config)
	if err != nil {
		return
	}

	if successFunc != nil {
		successFunc()
	}

	go ssh.DiscardRequests(reqs)
	go func() {
		for ch := range chans {
			ch.Reject(ssh.Prohibited, "Tunnels disallowed")
		}
	}()

	repeater = goutils.NewRepeater(time.Second*30, func() {
		if sshConn != nil {
			_, _, err := sshConn.SendRequest("ping", true, nil)
			if err != nil {
				log.Println(err)
			}
		}
	})

	channel := &Channel{
		conn: sshConn,
		comp: comp,
	}
	backend = ladder.NewBackEnd(channel)

	channels.AddBackEnd(backend)

	err = sshConn.Wait()

	channels.DelBackEnd(backend)

	repeater.Stop()
}
