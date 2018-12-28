package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

func generateServerConfig(authFunc func(ssh.ConnMetadata, []byte) (*ssh.Permissions, error)) (config *ssh.ServerConfig, err error) {
	var (
		key []byte
	)

	config = &ssh.ServerConfig{
		PasswordCallback: authFunc,
	}

	key, err = generateKey()
	if err != nil {
		return
	}
	private, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return
	}

	config.AddHostKey(private)
	return
}

func generateKey() (key []byte, err error) {
	var (
		r    = rand.Reader
		priv *ecdsa.PrivateKey
		bs   []byte
	)

	priv, err = ecdsa.GenerateKey(elliptic.P256(), r)
	if err != nil {
		return
	}

	bs, err = x509.MarshalECPrivateKey(priv)
	if err != nil {
		return
	}

	key = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: bs})
	return
}
