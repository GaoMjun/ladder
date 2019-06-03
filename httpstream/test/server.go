package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/GaoMjun/ladder"
	"github.com/GaoMjun/ladder/httpstream"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	var (
		upgrader     = httpstream.NewUpgrader()
		handleStream = func(stream *httpstream.Conn) {
			ladder.Pipe(stream, stream)
		}
	)

	go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bs, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(bs))
		fmt.Println("///////////////////")

		if len(r.Header.Get("Httpstream-Key")) > 0 {
			upgrader.Upgrade(w, r)
		}
	}))

	for {
		stream := upgrader.Accept()
		go handleStream(stream)
	}
}
