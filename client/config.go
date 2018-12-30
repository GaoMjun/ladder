package client

import (
	"encoding/json"
	"io/ioutil"
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
	confBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(confBytes, &config)
	return
}

func (self Config) String() (s string) {
	bs, _ := json.MarshalIndent(&self, "", "  ")
	s = string(bs)
	return
}
