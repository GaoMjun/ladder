package ladder

import (
	"io"
)

func Pipe(src io.ReadWriteCloser, dst io.ReadWriteCloser) {
	go func() {
		buffer := make([]byte, 1024)
		io.CopyBuffer(src, dst, buffer)
	}()

	buffer := make([]byte, 1024)
	io.CopyBuffer(dst, src, buffer)
}
