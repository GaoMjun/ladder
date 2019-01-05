package server

import (
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GaoMjun/ladder"
	"github.com/gorilla/websocket"
)

type server struct {
	listen string
	user   string
	pass   string
	key    [md5.Size]byte
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
		port := os.Getenv("PORT")
		if len(port) <= 0 {
			err = errors.New("invalid parameter")
			return
		}

		*l = ":" + port
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
	s.key = md5.Sum([]byte(fmt.Sprintf("%s:%s", s.user, s.pass)))

	err = http.ListenAndServe(s.listen, http.HandlerFunc(s.handler))
}

func (self *server) handler(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		token    string
		tokenOk  bool
		conn     *websocket.Conn
		upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
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

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	handleConn(ladder.NewConnWithXor(ladder.NewConn(conn), self.key[:]), self.user, self.pass)
}

func handleFake(w http.ResponseWriter, r *http.Request) {
	w.Write(index)
}
