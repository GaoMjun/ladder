package httpstream

import (
	"io"
	"net/http"
	"sync"
)

type Upgrader struct {
	readStreams  map[string]io.ReadCloser
	writeStreams map[string]io.WriteCloser
	locker       *sync.RWMutex
	connCh       chan *Conn
}

func NewUpgrader() (u *Upgrader) {
	u = &Upgrader{}
	u.readStreams = map[string]io.ReadCloser{}
	u.writeStreams = map[string]io.WriteCloser{}
	u.locker = &sync.RWMutex{}
	u.connCh = make(chan *Conn)
	return
}

func (self *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		return
	}

	key := r.Header.Get("Httpstream-Key")
	if len(key) <= 0 {
		return
	}

	if r.Method == "POST" {
		ws := self.getWriteStream(key)
		if ws == nil {
			self.addReadStream(key, r.Body)
		} else {
			self.connCh <- &Conn{r: r.Body, w: ws}
		}
	}

	if r.Method == "GET" {
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("X-Accel-Buffering", "no")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()

		rs := self.getReadStream(key)
		if rs == nil {
			self.addWriteStream(key, NopHttpResponseWriteCloser(w))
		} else {
			self.connCh <- &Conn{r: rs, w: NopHttpResponseWriteCloser(w)}
		}
	}

	<-w.(http.CloseNotifier).CloseNotify()
}

func (self *Upgrader) Accept() (conn *Conn) {
	conn = <-self.connCh
	return
}

func (self *Upgrader) getReadStream(key string) (r io.ReadCloser) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if s, ok := self.readStreams[key]; ok {
		r = s
		delete(self.readStreams, key)
		return
	}

	return
}

func (self *Upgrader) addReadStream(key string, r io.ReadCloser) {
	self.locker.Lock()
	defer self.locker.Unlock()

	self.readStreams[key] = r
	return
}

func (self *Upgrader) delReadStream(key string) {
	_ = self.getReadStream(key)
	return
}

func (self *Upgrader) getWriteStream(key string) (w io.WriteCloser) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if s, ok := self.writeStreams[key]; ok {
		w = s
		delete(self.writeStreams, key)
		return
	}

	return
}

func (self *Upgrader) addWriteStream(key string, w io.WriteCloser) {
	self.locker.Lock()
	defer self.locker.Unlock()

	self.writeStreams[key] = w
	return
}

func (self *Upgrader) delWriteStream(key string) {
	_ = self.getWriteStream(key)
	return
}
