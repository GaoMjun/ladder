package client

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/GaoMjun/ladder"
	"github.com/gorilla/websocket"
)

func Run(args []string) {
	var (
		err   error
		flags = flag.NewFlagSet("client", flag.ContinueOnError)

		config Config

		tcpServer *TCPServer
		channels  = ladder.NewChannels()
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	configFile := flags.String("c", "", "configuration file")
	flags.Parse(args)

	if len(*configFile) <= 0 {
		err = errors.New("invalid parameter")
		return
	}

	config, err = NewConfig(*configFile)
	if err != nil {
		return
	}
	fmt.Println(config)

	for _, remote := range config.Remotes {
		for i := 0; i < remote.Channels; i++ {
			go createChannel(config, remote, channels)
		}
	}

	tcpServer = NewTCPServer(config.Listen, channels)
	err = tcpServer.Run()
}

func createChannel(config Config, remote Remote, channels *ladder.Channels) {
	var (
		err       error
		user      = config.User
		pass      = config.Pass
		conn      *websocket.Conn
		token     string
		header    = map[string][]string{}
		dialer    = &websocket.Dialer{HandshakeTimeout: time.Second * 5}
		urlString = remote.Host
		u         *url.URL
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	u, err = url.Parse(urlString)
	if err != nil {
		return
	}

	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		err = errors.New(fmt.Sprint("not support protocol", u.Scheme))
		return
	}
	urlString = u.String()

	if len(remote.IP) > 0 {
		dialer.NetDial = func(network, addr string) (net.Conn, error) {
			return net.Dial(network, remote.IP)
		}
	}

	header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.3"}

TRY:
	token, _ = ladder.GenerateToken(user, pass)
	header["token"] = []string{token}

	if conn != nil {
		conn.Close()
		conn = nil
	}
	conn, _, err = dialer.Dial(urlString, header)
	if err != nil {
		log.Println(err)
		time.Sleep(time.Second * 3)
		goto TRY
	}

	log.Println("websocket connected")
	handleConn(config, ladder.NewConnWithSnappy(ladder.NewConn(conn)), channels)
	time.Sleep(time.Second * 3)
	goto TRY
}
