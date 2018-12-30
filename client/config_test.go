package client

import (
	"log"
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
