package ladder

import (
	socks5 "github.com/armon/go-socks5"
)

func NewSocks5Server() (server *socks5.Server) {
	conf := &socks5.Config{}
	server, _ = socks5.New(conf)

	return
}
