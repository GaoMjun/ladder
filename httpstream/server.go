package httpstream

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"sync"
)

type Upgrader struct {
	conns  map[string]*Conn
	locker *sync.RWMutex
	connCh chan *Conn
}

func NewUpgrader() (u *Upgrader) {
	u = &Upgrader{}
	u.conns = map[string]*Conn{}
	u.locker = &sync.RWMutex{}
	u.connCh = make(chan *Conn)
	return
}

func (self *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Httpstream-Key")
	remoteHost, _ := base64.StdEncoding.DecodeString(r.Header.Get("HTTPStream-Host"))

	if r.Method == "POST" {
		conn := self.getConn(key)
		if conn == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)

		data, _ := ioutil.ReadAll(r.Body)
		conn.dataCh <- data
		return
	}

	if r.Method == "GET" {
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()

		conn := &Conn{}
		conn.isClient = false
		conn.chunkedWriter = NopHttpResponseWriteCloser(w)
		conn.dataCh = make(chan []byte)
		conn.closeCh = make(chan struct{})
		conn.RemoteHost = string(remoteHost)

		self.addConn(key, conn)
		self.connCh <- conn

		<-w.(http.CloseNotifier).CloseNotify()
		self.delConn(key)

		underlyingConn, _, _ := w.(http.Hijacker).Hijack()
		underlyingConn.Close()
	}
}

func (self *Upgrader) Accept() (conn *Conn) {
	conn = <-self.connCh
	return
}

func (self *Upgrader) delConn(key string) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if c, ok := self.conns[key]; ok {
		close(c.closeCh)
		delete(self.conns, key)
		return
	}
	return
}

func (self *Upgrader) getConn(key string) (conn *Conn) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if c, ok := self.conns[key]; ok {
		conn = c
		return
	}

	return
}

func (self *Upgrader) addConn(key string, conn *Conn) {
	self.locker.Lock()
	defer self.locker.Unlock()

	if c, ok := self.conns[key]; ok {
		c.Close()
	}

	self.conns[key] = conn
	return
}
