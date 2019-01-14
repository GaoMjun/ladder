package httpstream

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type FormDataWriter struct {
	boundary string
	bufw     *bufio.Writer
}

func NewFormDataWriter() (fw *FormDataWriter) {
	fw = &FormDataWriter{}
	fw.boundary = randomBoundary()
	return
}

func (self *FormDataWriter) SetWriter(w io.Writer) {
	self.bufw = bufio.NewWriter(w)
}

func (self *FormDataWriter) Write(p []byte) (n int, err error) {
	if self.bufw == nil {
		err = errors.New("bufw is nil")
		return
	}

	fmt.Fprint(self.bufw, fmt.Sprint("--", self.boundary, "\r\n"))
	fmt.Fprint(self.bufw, "Content-Disposition: form-data; name=\"upload\"; filename=\"blob\"\r\n")
	fmt.Fprint(self.bufw, "Content-Type: application/octet-stream\r\n")
	fmt.Fprint(self.bufw, "Content-Transfer-Encoding: binary\r\n")
	fmt.Fprint(self.bufw, "\r\n")

	self.bufw.Write(p)
	fmt.Fprint(self.bufw, "\r\n")

	err = self.bufw.Flush()
	return
}

func (self *FormDataWriter) Close() (err error) {
	fmt.Fprint(self.bufw, fmt.Sprint("--", self.boundary, "--"))
	err = self.bufw.Flush()
	return
}

func (self *FormDataWriter) FormDataContentType() (s string) {
	s = "multipart/form-data; boundary=" + self.boundary
	return
}
