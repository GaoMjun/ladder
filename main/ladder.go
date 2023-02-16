package main

import (
	"errors"
	"log"
	"os"

	// _ "net/http/pprof"

	"github.com/GaoMjun/ladder/client"
	"github.com/GaoMjun/ladder/server"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe(":6061", nil))
	// }()

	var (
		err error
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "server")
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
