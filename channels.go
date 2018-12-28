package ladder

import (
	"container/list"
	"errors"
	"io"
	"sync"
)

type BackEnd struct {
	v interface{}
	e *list.Element
}

func NewBackEnd(v interface{}) (be *BackEnd) {
	be = &BackEnd{}
	be.v = v
	return
}

type Channels struct {
	backends *list.List
	index    uint
	locker   *sync.RWMutex
}

func NewChannels() (cs *Channels) {
	cs = &Channels{}
	cs.backends = list.New()
	cs.locker = &sync.RWMutex{}
	return
}

func (self *Channels) CreateStream() (stream io.ReadWriteCloser, err error) {
	self.locker.Lock()
	defer self.locker.Unlock()

	e := self.backends.Back()
	if e != nil {
		err = errors.New("no backend")
		return
	}
	stream = e.Value.(io.ReadWriteCloser)
	self.backends.MoveToFront(e)
	return
}

func (self *Channels) AddBackEnd(be *BackEnd) {
	self.locker.Lock()
	defer self.locker.Unlock()

	be.e = self.backends.PushBack(be)
}

func (self *Channels) DelBackEnd(be *BackEnd) {
	self.locker.Lock()
	defer self.locker.Unlock()

	self.backends.Remove(be.e)
}
