package client

import (
	"log"
	"net"
	"testing"
)

func TestConfig(t *testing.T) {
	config, err := NewConfig("config.jsonc")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(config)
}

func TestLookupHost(t *testing.T) {
	addrs, err := net.LookupHost("baidu.com")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(addrs)
}
