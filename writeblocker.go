package ladder

import (
	"io"
)

type WriteBlocker struct {
	w   io.Writer
	end chan bool
}

func NewWriteBlocker(w io.Writer, notify <-chan bool) (wb *WriteBlocker) {
	wb = &WriteBlocker{}
	wb.w = w
	wb.end = make(chan bool)
	go func() {
		<-notify
		wb.end <- true
	}()
	return
}

func (self *WriteBlocker) Write(p []byte) (n int, err error) {
	n, err = self.w.Write(p)
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
