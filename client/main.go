package client

import (
	"crypto/md5"
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

		channels = ladder.NewChannels()
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
			go createChannel(remote, channels)
		}
	}

	httpProxyServer := NewTCPServer("http", config.HttpListen, channels)
	socksProxyServer := NewTCPServer("socks", config.SocksListen, channels)

	go func() {
		err := httpProxyServer.Run()
		log.Panicln(err)
	}()

	err = socksProxyServer.Run()
}

func createChannel(remote Remote, channels *ladder.Channels) {
	var (
		err                error
		user               = remote.User
		pass               = remote.Pass
		conn               *websocket.Conn
		token              string
		header             = map[string][]string{}
		dialer             = &websocket.Dialer{HandshakeTimeout: time.Second * 5}
		urlString          = remote.Host
		u                  *url.URL
		dialFailedCount    = 0
		connectFailedCount = 0
		key                = md5.Sum([]byte(fmt.Sprintf("%s:%s", user, pass)))
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
		dialFailedCount++

		if dialFailedCount > 3 {
			return
		}
		time.Sleep(time.Second * 3)
		goto TRY
	}
	log.Println("websocket connected")

	connectFailedCount++
	if connectFailedCount > 3 {
		return
	}
	handleConn(user, pass, ladder.NewConnWithXor(ladder.NewConn(conn), key[:]), channels)
	time.Sleep(time.Second * 3)
	goto TRY
}
