package ladder

import (
	"fmt"
	"log"
	"testing"
)

func TestToken(t *testing.T) {
	var (
		err   error
		user  = "fuck"
		pass  = "gfw"
		token string
		ok    bool
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	token, err = GenerateToken(user, pass)
	if err != nil {
		return
	}

	log.Println(token)

	user, err = ParseToken(token, pass)
	if err != nil {
		return
	}

	log.Println(user)

	pass = "gfw"
	ok, err = CheckToken(user, pass, token)
	if err != nil {
		return
	}
	log.Println(ok)
}

func TestHeader(t *testing.T) {
	user := "fuck"
	pass := "gfw"
	iheader := "POST /api HTTP/1.1\r\nHost: 127.0.0.1:8000\r\nContent-type: application/json\r\nCache-Control: no-cache\r\nPostman-Token: 579eb349-e9da-45cf-dd38-5f1880140d65\r\n\r\n"

	o, err := EncryptHeader(iheader, user, pass)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(o)

	oheader, err := DecryptHeader(o, user, pass)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(oheader)
}
