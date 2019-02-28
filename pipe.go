package ladder

import (
	"io"
)

func Pipe(src, dst io.ReadWriteCloser) {
	go func() {
		var (
			err error
			buf = make([]byte, 1024*1)
			n   = 0
		)

		for {
			n, err = src.Read(buf)
			if err != nil {
				break
			}

			_, err = dst.Write(buf[:n])
			if err != nil {
				break
			}
		}
	}()

	var (
		err error
		buf = make([]byte, 1024*1)
		n   = 0
	)

	for {
		n, err = dst.Read(buf)
		if err != nil {
			break
		}

		_, err = src.Write(buf[:n])
		if err != nil {
			break
		}
	}

	src.Close()
	dst.Close()
}
