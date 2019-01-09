package client

import (
	"io"
	"log"

	"github.com/GaoMjun/ladder"
	"github.com/GaoMjun/ladder/mux"
)

func handleConn(host, user, pass string, comp bool, rc io.ReadCloser, channels *ladder.Channels, streamManager *ladder.StreamManager) {
	var (
		err     error
		backend *ladder.BackEnd
		stream  = mux.NewStream(rc, nil)
		frame   mux.Frame
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	channel := &Channel{
		host: host,
		user: user,
		pass: pass,
		comp: comp,
	}
	backend = ladder.NewBackEnd(channel)

	channels.AddBackEnd(backend)

	for {
		frame, err = stream.ReadFrame()
		if err != nil {
			break
		}

		rs := streamManager.GetReceiveStream(frame.StreamID)
		if rs != nil {
			rs.Ch <- frame.Data
		}
	}

	channels.DelBackEnd(backend)
}
