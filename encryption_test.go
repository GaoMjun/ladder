package ladder

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestEncrypt(t *testing.T) {
	var (
		err       error
		plantText = "hello"
		key       = []byte("haha000000000000")

		secret     []byte
		secretText string
		plan       []byte
	)
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	secret, err = encrypt([]byte(plantText), key)
	if err != nil {
		return
	}

	secretText = string(secret)
	log.Println(secretText)

	plan, err = decrypt(secret, key)
	if err != nil {
		return
	}
	plantText = string(plan)
	log.Println(plantText)
}
