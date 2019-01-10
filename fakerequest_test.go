package ladder

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"testing"
)

func TestFakeRequest(t *testing.T) {
	fakeRequest := NewFakeRequest("140.205.220.96:80", "token", "originHeader", "streamID", 0)

	conn, err := fakeRequest.Do()
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	response, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		log.Println(err)
		return
	}

	bs, err := httputil.DumpResponse(response, true)

	fmt.Println(string(bs))
}
