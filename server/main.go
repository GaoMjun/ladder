package server

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/GaoMjun/ladder"
	"github.com/gorilla/websocket"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type server struct {
	listen string
	user   string
	pass   string
}

func Run(args []string) {
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
	flags.Parse(args)

	if len(*l) <= 0 {
		err = errors.New("invalid parameter")
		return
	}

	if len(*u) <= 0 {
		err = errors.New("invalid parameter")
		return
	}

	if len(*p) <= 0 {
		err = errors.New("invalid parameter")
		return
	}

	s.listen = *l
	s.user = *u
	s.pass = *p

	err = http.ListenAndServe(s.listen, http.HandlerFunc(s.handler))
}

func (self *server) handler(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		token    string
		tokenOk  bool
		conn     *websocket.Conn
		upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	)
	defer func() {
		if err != nil {
			handleFake(w, r)
			log.Println(err)
		}
	}()

	token = r.Header.Get("token")
	if len(token) <= 0 {
		err = errors.New("token invalid")
		return
	}

	tokenOk, err = ladder.CheckToken(self.user, self.pass, token)
	if err != nil {
		return
	}

	if tokenOk != true {
		err = errors.New("token invalid")
		return
	}

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	handleConn(ladder.NewConnWithSnappy(ladder.NewConn(conn)), self.user, self.pass)
}

func handleFake(w http.ResponseWriter, r *http.Request) {
	w.Write(index)
}
