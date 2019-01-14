package ladderclient

import (
	"log"

	"github.com/GaoMjun/ladder/client"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func Run(configJsonString string) {
	client.RunWithJsonString(configJsonString)
}
