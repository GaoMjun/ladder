package ladderclient

import (
	"log"

	"github.com/GaoMjun/ladder/client"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func Run(filename string) {
	client.Run([]string{"-c", filename})
}
