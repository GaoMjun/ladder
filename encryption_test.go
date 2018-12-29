package ladder

import (
	"log"
	"testing"
)

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

	secret, err = Encrypt([]byte(plantText), key)
	if err != nil {
		return
	}

	secretText = string(secret)
	log.Println(secretText)

	plan, err = Decrypt(secret, key)
	if err != nil {
		return
	}
	plantText = string(plan)
	log.Println(plantText)
}
