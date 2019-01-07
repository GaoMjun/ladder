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
)

type server struct {
	listen   string
	user     string
	pass     string
	compress bool
	key      [md5.Size]byte
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
	m := flags.Bool("m", false, "compress message")
	flags.Parse(args)

	if len(*l) <= 0 {
		port := os.Getenv("PORT")
		if len(port) <= 0 {
			err = errors.New("invalid parameter, no listen address")
			return
		}

		*l = ":" + port
	}

	if len(*u) <= 0 {
		err = errors.New("invalid parameter, no user")
		return
	}

	if len(*p) <= 0 {
		err = errors.New("invalid parameter, no password")
		return
	}

	s.listen = *l
	s.user = *u
	s.pass = *p
	s.compress = *m
	s.key = md5.Sum([]byte(fmt.Sprintf("%s:%s", s.user, s.pass)))

	err = http.ListenAndServe(s.listen, http.HandlerFunc(s.handler))
}

func (self *server) handler(w http.ResponseWriter, r *http.Request) {
	var (
		err          error
		token        string
		tokenOk      bool
		originHeader string
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

	originHeader = r.Header.Get("Header")
	originHeader, err = ladder.DecryptHeader(originHeader, self.user, self.pass)
	if err != nil {
		return
	}

	handleConn(ladder.NewConnWithXor(ladder.NewConn(conn), self.key[:]), self.user, self.pass, self.compress)
}

func handleFake(w http.ResponseWriter, r *http.Request) {
	w.Write(index)
}
