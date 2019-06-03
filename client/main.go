// +build !lib

package client

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/GaoMjun/ladder/httpstream"

	"github.com/GaoMjun/ladder"
	"github.com/gorilla/websocket"
)

func RunWithJsonString(jsonString string) {
	var (
		err    error
		config Config
	)

	config, err = NewConfigWithJsonString(jsonString)
	if err != nil {
		return
	}

	err = run(config)
}

func Run(args []string) {
	var (
		err    error
		flags  = flag.NewFlagSet("client", flag.ContinueOnError)
		config Config
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

	err = run(config)
}

func run(config Config) (err error) {
	fmt.Println(config)

	var (
		channels = ladder.NewChannels()
	)

	for _, remote := range config.Remotes {
		for i := 0; i < remote.Channels; i++ {
			switch remote.Mode {
			case "ws":
				go createWSChannel(remote, channels)
			case "hs":
				go createHSChannel(remote, channels)
			default:
				go createWSChannel(remote, channels)
			}
		}
	}

	if len(config.HttpListen) > 0 {
		httpProxyServer := NewTCPServer("http", config.HttpListen, channels)
		go func() {
			log.Println(httpProxyServer.Run())
		}()

	}

	if len(config.SocksListen) > 0 {
		socksProxyServer := NewTCPServer("socks", config.SocksListen, channels)
		go func() {
			log.Println(socksProxyServer.Run())
		}()
	}

	if len(config.IPTransparentListen) > 0 {
		iptransparentProxyServer := NewTCPServer("iptransparent", config.IPTransparentListen, channels)
		go func() {
			log.Println(iptransparentProxyServer.Run())
		}()
	}

	select {}
}

func createWSChannel(remote Remote, channels *ladder.Channels) {
	var (
		err    error
		user   = remote.User
		pass   = remote.Pass
		comp   = remote.Compress
		conn   *websocket.Conn
		token  string
		header = map[string][]string{}
		dialer = &websocket.Dialer{HandshakeTimeout: time.Second * 5,
			ReadBufferSize: 1024, WriteBufferSize: 1024,
			TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
		urlString = remote.Host
		u         *url.URL
		key       = md5.Sum([]byte(fmt.Sprintf("%s:%s", user, pass)))
		reconnect = 0
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
	header["Token"] = []string{token}

	if conn != nil {
		conn.Close()
		conn = nil
	}
	conn, _, err = dialer.Dial(urlString, header)
	if err != nil {
		log.Println(err)

		reconnect = reconnectDuration(reconnect)
		time.Sleep(time.Second * time.Duration(reconnect))
		goto TRY
	}
	log.Println("connected to ", urlString)

	handleConn(user, pass, comp, ladder.NewConnWithXor(ladder.NewConn(conn), key[:]), channels, func() {
		reconnect = 1
	})

	reconnect = reconnectDuration(reconnect)
	time.Sleep(time.Second * time.Duration(reconnect))
	goto TRY
}

func createHSChannel(remote Remote, channels *ladder.Channels) {
	var (
		err       error
		dialer    = &httpstream.Dialer{}
		header    = http.Header{}
		token     string
		user      = remote.User
		pass      = remote.Pass
		comp      = remote.Compress
		key       = md5.Sum([]byte(fmt.Sprintf("%s:%s", user, pass)))
		conn      *httpstream.Conn
		reconnect = 0
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	if len(remote.IP) > 0 {
		dialer.NetDial = func(network, addr string) (net.Conn, error) {
			return net.Dial(network, remote.IP)
		}
	}

	if len(remote.UpHost) > 0 {
		dialer.UpHost = remote.UpHost
	}

	if len(remote.UpIP) > 0 {
		dialer.UpNetDial = func(network, addr string) (net.Conn, error) {
			return net.Dial(network, remote.UpIP)
		}
	}

	header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.3"}

TRY:
	token, _ = ladder.GenerateToken(user, pass)
	header["Token"] = []string{token}

	if conn != nil {
		conn.Close()
		conn = nil
	}
	conn, err = dialer.Dial(remote.Host, header)
	if err != nil {
		log.Println(err)

		reconnect = reconnectDuration(reconnect)
		time.Sleep(time.Second * time.Duration(reconnect))
		goto TRY
	}

	handleConn(user, pass, comp, ladder.NewConnWithXor(conn, key[:]), channels, func() {
		reconnect = 1
		log.Println("connected to ", remote.Host)
	})

	reconnect = reconnectDuration(reconnect)
	time.Sleep(time.Second * time.Duration(reconnect))
	goto TRY
}

func reconnectDuration(d1 int) (d2 int) {
	if d1 < 1 {
		d2 = 1
		return
	}

	d2 = d1 << 1
	if d2 > 60*3 {
		d2 = 1
	}
	return
}
