package ladder

import (
	"io"
)

type WriteBlocker struct {
	W   io.Writer
	end chan bool
}

func NewWriteBlocker(w io.Writer, notify <-chan bool) (wb *WriteBlocker) {
	wb = &WriteBlocker{}
	wb.W = w
	wb.end = make(chan bool)
	go func() {
		<-notify
		wb.end <- true
	}()
	return
}

func (self *WriteBlocker) Write(p []byte) (n int, err error) {
	n, err = self.W.Write(p)
	if err != nil {
		self.end <- true
	}
	return
}

func (self *WriteBlocker) Close() (err error) {
	return
}

func (self *WriteBlocker) Wait() {
	<-self.end
}
