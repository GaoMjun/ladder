package client

import (
	"io"
	"log"

	"github.com/GaoMjun/ladder"
	"github.com/GaoMjun/ladder/mux"
)

func handleConn(comp bool, rc io.ReadCloser, channels *ladder.Channels) {
	var (
		err     error
		backend *ladder.BackEnd
		stream  = mux.NewStream(rc, nil)
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	channel := &Channel{
		rc:   rc,
		comp: comp,
	}
	backend = ladder.NewBackEnd(channel) 

	channels.AddBackEnd(backend)

	for {
		frame, err = stream.ReadFrame()
		if err != nil {
			break
		}


	}

	channels.DelBackEnd(backend)
}
