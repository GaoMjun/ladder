package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	RawLines    []string
	HttpRequest *http.Request

	rawString string
	rawBytes  []byte
}

func NewRequest(r io.Reader) (request *Request, err error) {
	var (
		reader   = bufio.NewReader(r)
		line     []byte
		isPrefix bool
	)

	request = &Request{}

	for {
		line, isPrefix, err = reader.ReadLine()
		if err != nil {
			return
		}

		if isPrefix {
			err = errors.New("line is too long")
			return
		}

		request.RawLines = append(request.RawLines, string(line))

		if len(line) == 0 {
			break
		}
	}

	err = request.Parse()
	if err != nil {
		return
	}

	return
}

func (self *Request) Dump() (s string) {
	if len(self.rawString) > 0 {
		s = self.rawString
		return
	}

	s = strings.Join(self.RawLines, "\r\n")
	s = strings.Join([]string{s, "\r\n"}, "")

	self.rawString = s
	return
}

func (self *Request) DumpHex() (s string) {
	s = hex.Dump([]byte(self.Dump()))
	return
}

func (self *Request) Bytes() (bs []byte) {
	if len(self.rawBytes) > 0 {
		bs = self.rawBytes
		return
	}

	bs = []byte(self.Dump())

	self.rawBytes = bs
	return
}

func (self *Request) Parse() (err error) {
	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(self.Dump())))
	if err != nil {
		return
	}

	self.HttpRequest = r
	return
}
