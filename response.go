package ladder

import (
	"bufio"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	RawLines     []string
	HttpResponse *http.Response

	rawString string
	rawBytes  []byte
}

func NewResponse(r io.Reader) (response *Response, err error) {
	var (
		reader   = bufio.NewReader(r)
		line     []byte
		isPrefix bool
	)

	response = &Response{}

	for {
		line, isPrefix, err = reader.ReadLine()
		if err != nil {
			return
		}

		if isPrefix {
			err = errors.New("line is too long")
			return
		}

		response.RawLines = append(response.RawLines, string(line))

		if len(line) == 0 {
			break
		}
	}

	err = response.Parse()
	if err != nil {
		return
	}

	return
}

func (self *Response) String() (s string) {
	s = self.Dump()
	return
}

func (self *Response) Dump() (s string) {
	if len(self.rawString) > 0 {
		s = self.rawString
		return
	}

	s = strings.Join(self.RawLines, "\r\n")
	s = strings.Join([]string{s, "\r\n"}, "")

	self.rawString = s
	return
}

func (self *Response) DumpHex() (s string) {
	s = hex.Dump([]byte(self.Dump()))
	return
}

func (self *Response) Bytes() (bs []byte) {
	if len(self.rawBytes) > 0 {
		bs = self.rawBytes
		return
	}

	bs = []byte(self.Dump())

	self.rawBytes = bs
	return
}

func (self *Response) Parse() (err error) {
	r, err := http.ReadResponse(bufio.NewReader(strings.NewReader(self.Dump())), nil)
	if err != nil {
		return
	}

	self.HttpResponse = r
	return
}
