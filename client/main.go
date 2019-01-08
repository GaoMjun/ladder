package client

import (
	"bufio"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/GaoMjun/ladder"
)

func Run(args []string) {
	var (
		err   error
		flags = flag.NewFlagSet("client", flag.ContinueOnError)

		config Config

		channels = ladder.NewChannels()
		streamManager = ladder.NewStreamManager()
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

	httpProxyServer := NewTCPServer("http", config.HttpListen, channels, streamManager)
	socksProxyServer := NewTCPServer("socks", config.SocksListen, channels, streamManager)

	go func() {
		err := httpProxyServer.Run()
		log.Panicln(err)
	}()

	err = socksProxyServer.Run()
}

func createChannel(remote Remote, channels *ladder.Channels) {
	var (
		err         error
		user        = remote.User
		pass        = remote.Pass
		comp        = remote.Compress
		conn        net.Conn
		token       string
		urlString   = remote.Host
		u           *url.URL
		key         = md5.Sum([]byte(fmt.Sprintf("%s:%s", user, pass)))
		host        string
		fakeRequest *ladder.FakeRequest
		response    *http.Response
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

	if u.Scheme != "http" && u.Scheme != "https" {
		err = errors.New(fmt.Sprint("not support protocol", u.Scheme))
		return
	}

	switch u.Scheme {
	case "http":
		host = fmt.Sprint(u.Hostname(), "80")
	case "https":
		host = fmt.Sprint(u.Hostname(), "443")
	}

	if len(remote.IP) > 0 {
		host = remote.IP
	}

TRY:
	token, _ = ladder.GenerateToken(user, pass)
	fakeRequest = ladder.NewFakeRequest(host, token, "", "0")

	conn, err = fakeRequest.Do()
	if err != nil {
		log.Println(err)
		goto TRY
	}
	response, err = http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		log.Println(err)
		goto TRY
	}

	handleConn(comp, response.Body, channels)

	time.Sleep(time.Second * 3)
	goto TRY
}
