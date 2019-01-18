package server

import (
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	// _ "net/http/pprof"

	"github.com/GaoMjun/goutils"
	"github.com/GaoMjun/ladder"
	"github.com/GaoMjun/ladder/httpstream"
	"github.com/gorilla/websocket"
)

type server struct {
	listen     string
	user       string
	pass       string
	compress   bool
	key        [md5.Size]byte
	wsUpgrader websocket.Upgrader
	hsUpgrader *httpstream.Upgrader
}

func Run(args []string) {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6061", nil))
	// }()

	var (
		err   error
		flags = flag.NewFlagSet("server", flag.ContinueOnError)
		s     = &server{}
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	l := flags.String("l", "", "listen at")
	u := flags.String("u", "", "user")
	p := flags.String("p", "", "pass")
	m := flags.Bool("m", false, "compress message")
	flags.Parse(args)

	if len(*l) <= 0 {
		port := os.Getenv("PORT")
		if len(port) <= 0 {
			port = "80"
		}

		*l = ":" + port
	}

	if len(*u) <= 0 {
		*u = goutils.RandString(8)
	}

	if len(*p) <= 0 {
		*p = goutils.RandString(8)
	}

	log.Println("listen:", *l, "user:", *u, "pass:", *p)

	s.listen = *l
	s.user = *u
	s.pass = *p
	s.compress = *m
	s.key = md5.Sum([]byte(fmt.Sprintf("%s:%s", s.user, s.pass)))
	s.hsUpgrader = httpstream.NewUpgrader()
	s.wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	go func() {
		err = http.ListenAndServe(s.listen, http.HandlerFunc(s.handler))
		log.Panicln(err)
	}()

	for {
		conn := s.hsUpgrader.Accept()
		go s.handleHSConn(conn)
	}
}

func (self *server) handler(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		token   string
		tokenOk bool
		conn    *websocket.Conn
	)
	defer func() {
		if err != nil {
			handleFake(w, r)
			log.Println(err)
		}
	}()

	token = r.Header.Get("Token")
	if len(token) <= 0 {
		err = errors.New("token invalid, no token")
		return
	}

	tokenOk, err = ladder.CheckToken(self.user, self.pass, token)
	if err != nil {
		return
	}

	if tokenOk != true {
		err = errors.New(fmt.Sprint("token invalid,", token))
		return
	}

	hsKey := r.Header.Get("Httpstream-Key")
	if len(hsKey) > 0 {
		self.hsUpgrader.Upgrade(w, r)
		return
	}

	conn, err = self.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	handleConn(ladder.NewConnWithXor(ladder.NewConn(conn), self.key[:]), self.user, self.pass, self.compress)
}

func handleFake(w http.ResponseWriter, r *http.Request) {
	w.Write(index)
}

func (self *server) handleHSConn(conn *httpstream.Conn) {
	handleConn(ladder.NewConnWithXor(conn, self.key[:]), self.user, self.pass, self.compress)
}
