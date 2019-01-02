package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"regexp"
	"sync"
)

type Config struct {
	Listen string

	User string
	Pass string

	Remotes []Remote
}

type Remote struct {
	Host     string
	IP       string
	Channels int
}

func NewConfig(filename string) (config Config, err error) {
	var (
		confString string
		allRemote  = []Remote{}
		wg         = &sync.WaitGroup{}
		locker     = &sync.Mutex{}
	)

	confString, err = prepareConfig(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(confString), &config)
	if err != nil {
		return
	}

	for _, remote := range config.Remotes {
		wg.Add(1)
		go func(remote Remote) {
			defer wg.Done()

			rs := prepareRemote(remote)

			locker.Lock()
			allRemote = append(allRemote, rs...)
			locker.Unlock()

		}(remote)
	}

	wg.Wait()

	config.Remotes = allRemote
	return
}

func (self Config) String() (s string) {
	bs, _ := json.MarshalIndent(&self, "", "  ")
	s = string(bs)
	return
}

func prepareRemote(remote Remote) (remotes []Remote) {
	var (
		err   error
		u     *url.URL
		addrs []string
		port  string
	)
	defer func() {
		if len(remotes) <= 0 {
			remotes = append(remotes, remote)
		}

		if err != nil {
			log.Println(err)
		}
	}()

	if len(remote.IP) > 0 {
		return
	}

	u, err = url.Parse(remote.Host)
	if err != nil {
		return
	}

	port = u.Port()
	if len(port) <= 0 {
		switch u.Scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"
		default:
			err = errors.New(fmt.Sprint("not support protocol ", u.Scheme))
			return
		}
	}

	addrs, err = net.LookupHost(u.Hostname())
	if err != nil {
		return
	}

	for _, addr := range addrs {

		addr = fmt.Sprintf("%s:%s", addr, port)

		remote.IP = addr
		remotes = append(remotes, remote)
	}

	return
}

func prepareConfig(confPath string) (conf string, err error) {
	var (
		confBytes []byte
		re        *regexp.Regexp
	)

	confBytes, err = ioutil.ReadFile(confPath)
	if err != nil {
		return
	}

	conf = string(confBytes)

	re, err = regexp.Compile(`\s+//.*`)
	if err != nil {
		return
	}
	conf = re.ReplaceAllString(conf, "")

	re, err = regexp.Compile(`,\s*\n*}`)
	if err != nil {
		return
	}
	conf = re.ReplaceAllString(conf, "\n}")

	return
}
