package ladder

import (
	"errors"
	"sync"
)

type StreamManager struct {
	ids       map[uint16]uint16
	nousedids map[uint16]uint16
	start     uint16
	locker    *sync.RWMutex
}

type Stream struct {
	id uint16
	ch <-chan []byte
}

func NewStreamManager() (sm *StreamManager) {
	sm = &StreamManager{}
	sm.ids = map[uint16]uint16{}
	sm.nousedids = map[uint16]uint16{}
	sm.locker = &sync.RWMutex{}
	sm.start = 1
	return
}

func (self *StreamManager) GenerateStreamID() (id uint16, err error) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if len(self.nousedids) > 0 {
		for _, v := range self.nousedids {
			id = v
			break
		}

		delete(self.nousedids, id)
		self.ids[id] = id
		return
	}

	if self.start > ^uint16(0) {
		err = errors.New("no id available")
		return
	}

	id = self.start
	self.start++
	self.ids[id] = id
	return
}

func (self *StreamManager) RemoveStreamID(id uint16) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if _, ok := self.ids[id]; ok {
		delete(self.ids, id)
		self.nousedids[id] = id
	}
}
