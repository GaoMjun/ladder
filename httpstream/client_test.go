package httpstream

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/GaoMjun/ladder"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func TestStream(t *testing.T) {
	var (
		err      error
		conn     *Conn
		dialer   = &Dialer{}
		upgrader = NewUpgrader()
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	go http.ListenAndServe("127.0.0.1:8888", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader.Upgrade(w, r)
	}))

	time.Sleep(time.Second * 3)

	conn, err = dialer.Dial("http://127.0.0.1:8888/", nil)
	if err != nil {
		return
	}
	log.Println("connected")

	go func() {
		var (
			err    error
			n      int
			buffer = make([]byte, 1024)
		)
		defer func() {
			if err != nil {
				log.Println(err)
			}
		}()

		for {
			n, err = conn.Read(buffer)
			if err != nil {
				return
			}
			log.Println(string(buffer[:n]))
		}
	}()

	go func() {
		for {
			_, err = conn.Write([]byte("ping"))
			if err != nil {
				return
			}
			log.Println("ping to")

			time.Sleep(time.Second * 1)
		}
	}()

	for {
		stream := upgrader.Accept()
		go handleStream(stream)
	}
}

func handleStream(stream *Conn) {
	ladder.Pipe(stream, stream)
}
