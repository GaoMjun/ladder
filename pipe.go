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
			buf = make([]byte, 1024*4)
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
			buf = make([]byte, 1024*4)
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

func PipeIoCopy(src, dst io.ReadWriteCloser) {
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
		defer func() {
			o.Do(close)
			wg.Done()
		}()

		io.CopyBuffer(dst, src, make([]byte, 1024*4))
	}()

	go func() {
		defer func() {
			o.Do(close)
			wg.Done()
		}()

		io.CopyBuffer(src, dst, make([]byte, 1024*4))
	}()

	wg.Wait()
}
