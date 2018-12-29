package ladder

import (
	"container/list"
	"errors"
	"sync"
)

type BackEnd struct {
	V interface{}
	e *list.Element
}

func NewBackEnd(v interface{}) (be *BackEnd) {
	be = &BackEnd{}
	be.V = v
	return
}

type Channels struct {
	backends *list.List
	locker   *sync.RWMutex
}

func NewChannels() (cs *Channels) {
	cs = &Channels{}
	cs.backends = list.New()
	cs.locker = &sync.RWMutex{}
	return
}

func (self *Channels) GetBackEnd() (be *BackEnd, err error) {
	self.locker.Lock()
	defer self.locker.Unlock()

	e := self.backends.Back()
	if e == nil {
		err = errors.New("no backend")
		return
	}
	self.backends.MoveToFront(e)

	be = e.Value.(*BackEnd)
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
