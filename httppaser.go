package ladder

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

func ParseHttpHost(r io.Reader) (host string, port int, raw []byte, https bool, err error) {
	b := make([]byte, 1)
	if _, err = io.ReadFull(r, b); err != nil {
		return
	}

	if b[0] == 'G' ||
		b[0] == 'H' ||
		b[0] == 'P' ||
		b[0] == 'O' {
		if host, port, raw, err = parseHttp(r); err != nil {
			return
		}
		raw = append(b, raw...)
		return
	}

	if b[0] == 0x16 {
		if host, port, raw, err = parseHttps(r); err != nil {
			return
		}
		https = true
		return
	}

	raw = b
	err = errors.New("not http protocol")
	return
}

func parseHttp(conn io.Reader) (host string, port int, raw []byte, err error) {
	var (
		b  = make([]byte, 1024)
		n  = 0
		ss []string
	)

	if n, err = conn.Read(b); err != nil {
		return
	}
	raw = b[:n]

	ss = strings.Split(string(raw), "\r\n")
	if len(ss) <= 0 {
		return
	}
	if strings.HasPrefix(ss[0], "ET ") ||
		strings.HasPrefix(ss[0], "EAD ") ||
		strings.HasPrefix(ss[0], "OST ") ||
		strings.HasPrefix(ss[0], "PTIONS ") {

		for _, s := range ss[1:] {
			kv := strings.Split(s, ":")
			k := kv[0]
			if strings.ToLower(k) == "host" {
				host, port = splitHostPort(strings.TrimSpace(kv[1]))
			}
		}
	}
	return
}

func parseHttps(conn io.Reader) (host string, port int, raw []byte, err error) {
	raw = append(raw, byte(0x16))

	// Version: TLS 1.0
	b := make([]byte, 2)
	if _, err = io.ReadFull(conn, b); err != nil {
		return
	}
	raw = append(raw, b...)
	if !(b[0] == 0x03 && b[1] == 0x01) {
		err = errors.New("parseHttps: parse failed")
		return
	}

	// Length
	b = make([]byte, 2)
	if _, err = io.ReadFull(conn, b); err != nil {
		return
	}
	raw = append(raw, b...)
	length := int(b[0])<<8 | int(b[1])
	raw = make([]byte, 1+2+2+length)
	raw[0] = 0x16
	raw[1] = 0x03
	raw[2] = 0x01
	raw[3] = b[0]
	raw[4] = b[1]

	if _, err = io.ReadFull(conn, raw[5:]); err != nil {
		raw = raw[:5]
		return
	}

	// Handshake
	offset := 5

	// Handshake Type: Client Hello
	if raw[offset] != 0x01 {
		err = errors.New("parseHttps: parse failed")
		return
	}
	offset++
	if offset > len(raw)-1 {
		err = errors.New("parseHttps: parse failed")
		return
	}

	// Length
	length = int(raw[offset])<<16 | int(raw[offset+1])<<8 | int(raw[offset+2])
	offset += 3
	if offset > len(raw)-1 {
		err = errors.New("parseHttps: parse failed")
		return
	}

	// Version: TLS
	offset += 2
	if offset > len(raw)-1 {
		err = errors.New("parseHttps: parse failed")
		return
	}

	// Random (32 bytes fixed length)
	offset += 32
	if offset > len(raw)-1 {
		err = errors.New("parseHttps: parse failed")
		return
	}

	// Session ID Length
	length = int(raw[offset])
	offset++
	offset += length
	if offset > len(raw)-1 {
		err = errors.New("parseHttps: parse failed")
		return
	}

	// Cipher Suites Length
	length = int(raw[offset])<<8 | int(raw[offset+1])
	offset += 2
	offset += length
	if offset > len(raw)-1 {
		err = errors.New("parseHttps: parse failed")
		return
	}

	// Compression Methods Length
	length = int(raw[offset])
	offset++
	offset += length
	if offset > len(raw)-1 {
		err = errors.New("parseHttps: parse failed")
		return
	}

	// Extensions Length
	length = int(raw[offset])<<8 | int(raw[offset+1])
	offset += 2
	if offset > len(raw)-1 {
		err = errors.New("parseHttps: parse failed")
		return
	}

	for {
		if offset > len(raw)-1 {
			err = errors.New("parseHttps: parse failed")
			return
		}

		t := int(raw[offset])<<8 | int(raw[offset+1])
		offset += 2
		if offset > len(raw)-1 {
			err = errors.New("parseHttps: parse failed")
			return
		}

		length = int(raw[offset])<<8 | int(raw[offset+1])
		offset += 2
		if offset > len(raw)-1 {
			err = errors.New("parseHttps: parse failed")
			return
		}

		if t == 0 {
			// Server Name Indication
			length = int(raw[offset])<<8 | int(raw[offset+1])
			offset += 2
			if offset > len(raw)-1 {
				err = errors.New("parseHttps: parse failed")
				return
			}

			if raw[offset] != 0 {
				err = errors.New("parseHttps: parse failed")
				return
			}
			offset++
			if offset > len(raw)-1 {
				err = errors.New("parseHttps: parse failed")
				return
			}

			length = int(raw[offset])<<8 | int(raw[offset+1])
			offset += 2
			if offset+length > len(raw)-1 {
				err = errors.New("parseHttps: parse failed")
				return
			}

			host = string(raw[offset : offset+length])
			break
		}

		offset += length
	}

	port = 443
	return
}

func splitHostPort(hostport string) (host string, port int) {
	ss := strings.Split(hostport, ":")
	host = ss[0]

	if len(ss) == 2 {
		port, _ = strconv.Atoi(ss[1])
		return
	}

	port = 80
	return
}
