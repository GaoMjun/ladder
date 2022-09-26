package ladder

import (
	"crypto/md5"
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

func TestXor(t *testing.T) {
	i := []byte("hello world")
	key := md5.Sum([]byte("user:pass"))

	log.Println(i)
	log.Println(key)

	o := xor(i, key[:])
	log.Println(o)

	// o = xor(o, key[:])
	// log.Println(o)

	for _, b := range o {
		c := []byte{b}
		d := xor(c, key[:])
		log.Println(d)
	}
}
