package ladder

import (
	"errors"
	"sync"
)

type StreamManager struct {
	streamID       *StreamID
	receiveStreams map[uint16]*ReceiveStream
	locker         *sync.RWMutex
}

type StreamID struct {
	ids       map[uint16]uint16
	nousedids map[uint16]uint16
	start     uint16
	locker    *sync.RWMutex
}

func NewStreamID() (s *StreamID) {
	s = &StreamID{}
	s.ids = map[uint16]uint16{}
	s.nousedids = map[uint16]uint16{}
	s.locker = &sync.RWMutex{}
	s.start = 1
	return
}

func NewStreamManager() (sm *StreamManager) {
	sm = &StreamManager{}

	sm.streamID = NewStreamID()
	sm.receiveStreams = map[uint16]*ReceiveStream{}
	sm.locker = &sync.RWMutex{}
	return
}

func (self *StreamManager) GenerateStreamID() (id uint16, err error) {
	self.streamID.locker.Lock()
	defer self.streamID.locker.Unlock()

	if len(self.streamID.nousedids) > 0 {
		for _, v := range self.streamID.nousedids {
			id = v
			break
		}

		delete(self.streamID.nousedids, id)
		self.streamID.ids[id] = id
		return
	}

	if self.streamID.start > ^uint16(0) {
		err = errors.New("no id available")
		return
	}

	id = self.streamID.start
	self.streamID.start++
	self.streamID.ids[id] = id
	return
}

func (self *StreamManager) RemoveStreamID(id uint16) {
	self.streamID.locker.Lock()
	defer self.streamID.locker.Unlock()

	if _, ok := self.streamID.ids[id]; ok {
		delete(self.streamID.ids, id)
		self.streamID.nousedids[id] = id
	}
}

type ReceiveStream struct {
	id     uint16
	Ch     chan []byte
	buffer []byte
}

func (self *StreamManager) NewReceiveStream(id uint16) (rs *ReceiveStream) {
	rs = &ReceiveStream{}
	rs.id = id
	rs.Ch = make(chan []byte, 1024)
	self.AddReceiveStream(rs)
	return
}

func (self *StreamManager) AddReceiveStream(rs *ReceiveStream) {
	self.locker.Lock()
	defer self.locker.Unlock()

	self.receiveStreams[rs.id] = rs
}

func (self *StreamManager) DelReceiveStream(id uint16) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if _, ok := self.receiveStreams[id]; ok {
		delete(self.receiveStreams, id)
	}
}

func (self *StreamManager) GetReceiveStream(id uint16) (rs *ReceiveStream) {
	self.locker.RLock()
	defer self.locker.RUnlock()

	rs = self.receiveStreams[id]
	return
}

func (self *ReceiveStream) Read(buf []byte) (n int, err error) {
	if len(self.buffer) > 0 {
		if len(self.buffer) <= len(buf) {
			n = copy(buf, self.buffer)
			self.buffer = []byte{}
			return
		}

		n = copy(buf, self.buffer)
		self.buffer = self.buffer[n:]
		return
	}

	data := <-self.Ch

	if len(data) <= len(buf) {
		n = copy(buf, data)
		return
	}

	n = copy(buf, data)
	self.buffer = data[n:]
	return
}

func (self *ReceiveStream) Close() (err error) {
	return
}
