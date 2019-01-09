package mux

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

var SIGNATURE [3]byte = [3]byte{'m', 'u', 'x'}
var VERSION uint8 = 1

type Frame struct {
	Signature [3]byte
	Version   uint8
	StreamID  uint16
	Length    uint16
	Data      []byte
}

func (self Frame) String() (s string) {
	t := self
	t.Data = []byte{}

	bs, _ := json.MarshalIndent(&t, "", "  ")
	s = string(bs)
	return
}

func (self Frame) HeaderMarshal() (bs []byte) {
	buffer := &bytes.Buffer{}

	binary.Write(buffer, binary.LittleEndian, self.Signature)
	binary.Write(buffer, binary.BigEndian, self.Version)
	binary.Write(buffer, binary.BigEndian, self.StreamID)
	binary.Write(buffer, binary.BigEndian, self.Length)

	bs = buffer.Bytes()
	return
}
