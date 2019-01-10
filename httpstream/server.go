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

func (self *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.Method != "GET" {
		return
	}

	if len(r.Header["HTTPStream-Key"]) <= 0 {
		return
	}

	var (
		key = r.Header["HTTPStream-Key"][0]
	)

	if r.Method == "POST" {
		w := self.getWriteStream(key)
		if w == nil {
			self.addReadStream(key, r.Body)
		} else {
			self.connCh <- &Conn{r: r.Body, w: w}
		}
	}

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "octet-stream")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()

		r := self.getReadStream(key)
		if r == nil {
			self.addWriteStream(key, NopWriteCloser(w))
		} else {
			self.connCh <- &Conn{r: r, w: NopWriteCloser(w)}
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
