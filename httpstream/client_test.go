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
		err    error
		stream *HTTPStream
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	go http.ListenAndServe("127.0.0.1:8888", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	time.Sleep(time.Second * 3)

	dialer = &Dialer{}
	stream, err = OpenStream("http://127.0.0.1:8888/")
	if err != nil {
		return
	}

	// go func() {
	// 	var (
	// 		err    error
	// 		n      int
	// 		buffer = make([]byte, 1024)
	// 	)
	// 	defer func() {
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 	}()

	// 	for {
	// 		n, err = stream.Read(buffer)
	// 		if err != nil {
	// 			return
	// 		}
	// 		log.Println(string(buffer[:n]))
	// 	}
	// }()

	// for {
	// 	_, err = stream.Write([]byte("ping"))
	// 	if err != nil {
	// 		return
	// 	}
	// 	// log.Println("ping to")

	// 	time.Sleep(time.Second * 1)
	// }

	ladder.Pipe(stream, stream)
}
