package ladder

import (
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

	pass = "gfw1"
	ok, err = CheckToken(user, pass, token)
	log.Println(ok)
}
