package ladder

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func encrypt(i, key []byte) (o []byte, err error) {
	var (
		block cipher.Block
		mode  cipher.BlockMode
	)

	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}

	mode = cipher.NewCBCEncrypter(block, key)

	i, _ = padding(i, block.BlockSize())
	o = make([]byte, len(i))

	mode.CryptBlocks(o, i)
	return
}

func padding(i []byte, blockSize int) (o []byte, err error) {
	var (
		paddingSize int
		padding     []byte
	)

	paddingSize = blockSize - len(i)%blockSize
	padding = bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	o = append(i, padding...)
	return
}

func decrypt(i, key []byte) (o []byte, err error) {
	var (
		block cipher.Block
		mode  cipher.BlockMode
	)

	block, err = aes.NewCipher(key)
	if err != nil {
		return
	}

	o = make([]byte, len(i))
	mode = cipher.NewCBCDecrypter(block, key)
	mode.CryptBlocks(o, i)

	o, err = unPadding(o, block.BlockSize())
	return
}

func unPadding(i []byte, blockSize int) (o []byte, err error) {
	var (
		size        = len(i)
		paddingSize = int(i[size-1])
	)

	if paddingSize < 0 || paddingSize > blockSize {
		err = errors.New("uppadding failed")
		return
	}

	o = i[:(size - paddingSize)]
	return
}
