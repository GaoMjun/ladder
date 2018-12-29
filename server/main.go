package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/GaoMjun/ladder"
	"github.com/gorilla/websocket"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	log.Println(http.ListenAndServe(":8888", http.HandlerFunc(handler)))
}

func handler(w http.ResponseWriter, r *http.Request) {
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

	tokenOk, err = ladder.CheckToken("fuck", "gfw", token)
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

	handleConn(ladder.NewConnWithSnappy(ladder.NewConn(conn)))
}

func handleFake(w http.ResponseWriter, r *http.Request) {
	w.Write(index)
}
