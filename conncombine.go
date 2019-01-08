package ladder

import "io"

type ConnCombine struct {
	wc io.WriteCloser
	rc io.ReadCloser
}

func NewConnCombine(wc io.WriteCloser, rc io.ReadCloser) (c *ConnCombine) {
	c = &ConnCombine{}
	c.wc = wc
	c.rc = rc
	return
}

func (self *ConnCombine) Read(p []byte) (n int, err error) {
	n, err = self.rc.Read(p)
	return
}

func (self *ConnCombine) Write(p []byte) (n int, err error) {
	n, err = self.wc.Write(p)
	return
}

func (self *ConnCombine) Close() (err error) {
	err = self.wc.Close()
	if err != nil {
		return
	}

	err = self.rc.Close()
	if err != nil {
		return
	}
	return
}
