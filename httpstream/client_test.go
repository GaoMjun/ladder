package httpstream

import (
	"log"
	"net/http"
	"testing"
	"time"

	_ "net/http/pprof"

	"ladder"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	go func() {
		log.Println(http.ListenAndServe(":6061", nil))
	}()
}

func TestStream(t *testing.T) {
	var (
		err      error
		conn     *Conn
		upgrader = NewUpgrader()
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	go http.ListenAndServe("127.0.0.1:8888", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get("Httpstream-Key")) > 0 {
			upgrader.Upgrade(w, r)
		}
	}))

	// select {}
	time.Sleep(time.Second * 3)

	if conn, err = Dial("wotv.17wo.cn", "127.0.0.1:8888", nil); err != nil {
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

	time.AfterFunc(time.Second*5, func() {
		conn.Close()
	})
	go func() {
		for {
			_, err = conn.Write([]byte("pingping"))
			if err != nil {
				return
			}
			// log.Println("ping to")

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
