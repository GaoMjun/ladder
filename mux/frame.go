package mux

import (
	"bytes"
	"encoding/binary"
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

func (self Frame) HeaderMarshal() (bs []byte) {
	buffer := &bytes.Buffer{}

	binary.Write(buffer, binary.LittleEndian, self.Signature)
	binary.Write(buffer, binary.LittleEndian, self.Version)
	binary.Write(buffer, binary.BigEndian, self.StreamID)
	binary.Write(buffer, binary.LittleEndian, self.Length)

	bs = buffer.Bytes()
	return
}
