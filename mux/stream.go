package mux

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Stream struct {
	r   *bufio.Reader
	w   *bufio.Writer
	id  uint16
	buf []byte
}

func NewStream(rd io.Reader, w io.Writer) (stream *Stream) {
	stream = &Stream{}

	stream.r = bufio.NewReader(rd)
	stream.w = bufio.NewWriter(w)
	return
}

func (self *Stream) ReadFrame() (frame Frame, err error) {
	var (
		sig  = make([]byte, 3)
		ver  = make([]byte, 1)
		sid  = make([]byte, 2)
		len  = make([]byte, 2)
		data []byte
	)
	frame = Frame{}

	_, err = io.ReadFull(self.r, sig)
	if err != nil {
		return
	}
	if !bytes.Equal(sig, SIGNATURE[:]) {
		err = errors.New("read signature failed")
		return
	}
	frame.Signature = SIGNATURE

	_, err = io.ReadFull(self.r, ver)
	if err != nil {
		return
	}
	if uint8(ver[0]) != VERSION {
		err = errors.New("read version failed")
		return
	}
	frame.Version = VERSION

	_, err = io.ReadFull(self.r, sid)
	if err != nil {
		return
	}
	frame.StreamID = binary.BigEndian.Uint16(sid)

	_, err = io.ReadFull(self.r, len)
	if err != nil {
		return
	}
	frame.Length = binary.BigEndian.Uint16(len)

	data = make([]byte, frame.Length)
	_, err = io.ReadFull(self.r, data)
	if err != nil {
		return
	}

	frame.Data = data
	return
}

func (self *Stream) Read(data []byte) (n int, err error) {
	if len(self.buf) > 0 {
		if len(self.buf) <= len(data) {
			n = copy(data, self.buf)
			self.buf = []byte{}
			return
		}

		n = copy(data, self.buf)
		self.buf = self.buf[n:]
		return
	}

	frame, err := self.ReadFrame()
	if err != nil {
		return
	}

	if int(frame.Length) <= len(data) {
		n = copy(data, frame.Data)
		return
	}

	n = copy(data, frame.Data)
	self.buf = frame.Data[n:]
	return
}

func (self *Stream) WriteFrame(frame Frame) (err error) {
	_, err = self.w.Write(frame.HeaderMarshal())
	if err != nil {
		return
	}

	_, err = self.w.Write(frame.Data)
	return
}

func (self *Stream) Write(data []byte) (n int, err error) {
	var (
		frame  = Frame{Signature: SIGNATURE, Version: VERSION, StreamID: self.id}
		offset = 0
		length = 0
	)
	defer func() {
		if err == nil {
			n = len(data)
		}
	}()

	for {
		length = len(data) - offset
		if length <= 0 {
			return
		}

		if length <= int(^uint16(0)) {
			frame.Length = uint16(length)
			frame.Data = data[offset : offset+int(frame.Length)]
			err = self.WriteFrame(frame)
			if err != nil {
				return
			}

			offset += int(frame.Length)
			continue
		}

		frame.Length = ^uint16(0)
		frame.Data = data[offset : offset+int(frame.Length)]
		err = self.WriteFrame(frame)
		if err != nil {
			return
		}

		offset += int(frame.Length)
	}

	return
}

func (self *Stream) Flush() (err error) {
	err = self.w.Flush()
	return
}
