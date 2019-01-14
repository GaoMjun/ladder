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
	HttpListen  string
	SocksListen string
	Remotes     []Remote
}

type Remote struct {
	Host     string
	IP       string
	Channels int
	User     string
	Pass     string
	Compress bool
	Mode     string
	UpHost   string
	UpIP     string
}

func NewConfigWithJsonString(jsonString string) (config Config, err error) {
	var (
		confString string
		allRemote  = []Remote{}
		wg         = &sync.WaitGroup{}
		locker     = &sync.Mutex{}
	)

	confString, err = prepareConfig(jsonString)
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

func NewConfig(filename string) (config Config, err error) {
	var (
		confString string
		allRemote  = []Remote{}
		wg         = &sync.WaitGroup{}
		locker     = &sync.Mutex{}
		confBytes  []byte
		jsonString string
	)

	confBytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	jsonString = string(confBytes)

	confString, err = prepareConfig(jsonString)
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

func prepareConfig(jsonString string) (conf string, err error) {
	var (
		re *regexp.Regexp
	)

	conf = jsonString

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
