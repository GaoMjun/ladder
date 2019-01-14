package httpstream

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type FormDataReader struct {
	r    io.ReadCloser
	bufr *bufio.Reader
}

func NewFormDataReader(r io.ReadCloser) (fr *FormDataReader) {
	fr = &FormDataReader{}
	fr.r = r
	fr.bufr = bufio.NewReader(r)
	return
}

func (self *FormDataReader) Read(p []byte) (n int, err error) {
	err = self.readPart()
	return
}

func (self *FormDataReader) Close() (err error) {
	err = self.r.Close()
	return
}

func (self *FormDataReader) readPart() (err error) {
	var (
		line     []byte
		isPrefix bool
	)

	for {
		line, isPrefix, err = self.bufr.ReadLine()
		if err != nil {
			return
		}
		fmt.Println(string(line))

		if isPrefix {
			err = errors.New("line is too long")
			return
		}

		if len(line) == 0 {
			break
		}
	}

	return
}
