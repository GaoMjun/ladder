package ladder

import (
	"io"
	"sync"
)

func Pipe(src io.ReadWriteCloser, dst io.ReadWriteCloser) (int64, int64) {

	var sent, received int64
	var c = make(chan bool)
	var o sync.Once

	close := func() {
		src.Close()
		dst.Close()
		close(c)
	}

	go func() {
		received, _ = io.CopyBuffer(src, dst, make([]byte, 1024))
		o.Do(close)
	}()

	go func() {
		sent, _ = io.CopyBuffer(dst, src, make([]byte, 1024))
		o.Do(close)
	}()

	<-c
	return sent, received
}
