package ladder

import (
	"io"
	"sync"
)

func Pipe(src, dst io.ReadWriteCloser) {
	var (
		close = func() {
			src.Close()
			dst.Close()
		}
		wg = sync.WaitGroup{}
		o  = sync.Once{}
	)

	wg.Add(2)
	go func() {
		var (
			err error
			buf = make([]byte, 1024*1)
			n   = 0
		)
		defer func() {
			o.Do(close)
			wg.Done()
		}()

		for {
			if n, err = src.Read(buf); err != nil {
				return
			}

			if _, err = dst.Write(buf[:n]); err != nil {
				return
			}
		}
	}()

	go func() {
		var (
			err error
			buf = make([]byte, 1024*1)
			n   = 0
		)
		defer func() {
			o.Do(close)
			wg.Done()
		}()

		for {
			if n, err = dst.Read(buf); err != nil {
				return
			}

			if _, err = src.Write(buf[:n]); err != nil {
				return
			}
		}
	}()

	wg.Wait()
}
