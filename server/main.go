package main

import (
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

	token = r.Header.Get("token")
	log.Println(token)
	tokenOk, err = ladder.CheckToken("fuck", "gfw", token)
	if err != nil {
		return
	}
	log.Println(tokenOk)

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	handleConn(ladder.NewConn(conn))
}
