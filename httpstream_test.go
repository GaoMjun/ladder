package ladder

import (
	"log"
	"net/http"
	"testing"
	"time"
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
		if r.Method == "POST" {
			buffer := make([]byte, 1024)
			for {
				n, err := r.Body.Read(buffer)
				if err != nil {
					return
				}

				// log.Println(string(buffer[:n]))
				if string(buffer[:n]) == "ping" {
					log.Println("pong")
				}
			}
		}

		if r.Method == "GET" {
			w.Header().Set("Content-Type", "octet-stream")
			w.Header().Set("Transfer-Encoding", "chunked")
			w.WriteHeader(http.StatusOK)
			w.(http.Flusher).Flush()

			// wc := httputil.NewChunkedWriter(w)
			for {
				_, err := w.Write([]byte("ping"))
				if err != nil {
					return
				}
				w.(http.Flusher).Flush()
				log.Println("ping")

				time.Sleep(time.Second * 1)
			}
		}
	}))

	time.Sleep(time.Second * 3)

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

	Pipe(stream, stream)
}
