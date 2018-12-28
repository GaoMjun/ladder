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
		conn     *websocket.Conn
		upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	)

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	handleConn(ladder.NewConn(conn))
}
