package main

import (
	"errors"
	"log"
	"os"

	"github.com/GaoMjun/ladder/client"
	"github.com/GaoMjun/ladder/server"
)

func main() {
	var (
		err error
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	if len(os.Args) < 2 {
		err = errors.New("invalid parameter")
		return
	}

	mode := os.Args[1]

	if mode != "server" && mode != "client" {
		err = errors.New("invalid parameter")
		return
	}

	if mode == "server" {
		server.Run(os.Args[2:])
		return
	}

	if mode == "client" {
		client.Run(os.Args[2:])
		return
	}
}