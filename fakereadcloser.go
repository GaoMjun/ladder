package ladder

import "io"

type FakeReadCloser struct {
	end chan bool
}

func NewFakeReadCloser() (fr *FakeReadCloser) {
	fr = &FakeReadCloser{}
	fr.end = make(chan bool)
	return
}

func (self *FakeReadCloser) Read(p []byte) (n int, err error) {
	<-self.end
	err = io.EOF
	return
}

func (self *FakeReadCloser) Close() (err error) {
	self.end <- true
	return
}
