package ladder

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateToken(user, pass string) (token string, err error) {
	var (
		plant  = fmt.Sprintf("%v:%v", user, time.Now().Unix())
		key    = md5.Sum([]byte(pass))
		secret []byte
	)

	secret, err = encrypt([]byte(plant), key[:])
	if err != nil {
		return
	}

	token = base64.StdEncoding.EncodeToString(secret)
	return
}

func ParseToken(token string, pass string) (user string, err error) {
	var (
		key        = md5.Sum([]byte(pass))
		secret     []byte
		plantBytes []byte
		plan       string
		ss         []string
		timestamp  int64
		duration   int64
	)

	secret, err = base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}

	plantBytes, err = decrypt(secret, key[:])
	if err != nil {
		return
	}
	plan = string(plantBytes)

	ss = strings.Split(plan, ":")
	if len(ss) != 2 {
		err = errors.New("token invalid")
		return
	}

	timestamp, err = strconv.ParseInt(ss[1], 10, 64)
	if err != nil {
		err = errors.New("token invalid, timestamp invalid")
		return
	}

	duration = time.Now().Unix() - timestamp
	if duration < -60 || duration > 60 {
		err = errors.New(fmt.Sprint("token invalid, timeout ", duration))
		return
	}

	user = ss[0]

	return
}

func CheckToken(user, pass, token string) (ok bool, err error) {
	_user, err := ParseToken(token, pass)
	if err != nil {
		return
	}

	if _user != user {
		ok = false
		return
	}

	ok = true

	return
}
