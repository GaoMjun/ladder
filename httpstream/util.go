package httpstream

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func generateChallengeKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func NopWriteCloser(w io.Writer) (wc io.WriteCloser) {
	return nopWriteCloser{w}
}

type nopHttpResponseWriteCloser struct {
	w http.ResponseWriter
}

func (self nopHttpResponseWriteCloser) Write(p []byte) (n int, err error) {
	n, err = self.w.Write(p)
	if err != nil {
		return
	}
	self.w.(http.Flusher).Flush()
	return
}

func (self nopHttpResponseWriteCloser) Close() error {
	// if flusher, ok := self.w.(http.Flusher); ok {
	// 	flusher.Flush()
	// }
	return nil
}

func NopHttpResponseWriteCloser(w http.ResponseWriter) (wc io.WriteCloser) {
	return &nopHttpResponseWriteCloser{w}
}
