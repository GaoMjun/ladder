package main

import (
	"log"

	"github.com/GaoMjun/ladder"
	"github.com/gorilla/websocket"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var (
		err  error
		conn *websocket.Conn

		tcpServer *TCPServer
		channels  = ladder.NewChannels()
		token     string
		header    map[string][]string
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	token, _ = ladder.GenerateToken("fuck", "gfw")
	header = map[string][]string{"token": []string{token}}

	for i := 0; i < 10; i++ {
		go func() {
			conn, _, err = websocket.DefaultDialer.Dial("ws://192.168.1.57:8888/", header)
			if err != nil {
				return
			}

			handleConn(ladder.NewConn(conn), channels)
		}()
	}

	tcpServer = NewTCPServer(":9999", channels)
	err = tcpServer.Run()
}
