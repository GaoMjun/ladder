package client

import (
	"golang.org/x/crypto/ssh"
)

type Channel struct {
	conn ssh.Conn
	comp bool
}
