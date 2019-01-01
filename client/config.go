package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
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
		confBytes []byte
		allRemote = []Remote{}
		wg        = &sync.WaitGroup{}
		locker    = &sync.Mutex{}
	)

	confBytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(confBytes, &config)
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

	addrs, err = net.LookupHost(u.Hostname())
	if err != nil {
		return
	}

	for _, addr := range addrs {
		if u.Scheme == "https" {
			addr = fmt.Sprintf("%s:443", addr)
		} else {
			addr = fmt.Sprintf("%s:80", addr)
		}

		remote.IP = addr
		remotes = append(remotes, remote)
	}

	return
}
