module github.com/GaoMjun/ladder

go 1.19

replace github.com/GaoMjun/iptransparent => ../iptransparent

replace github.com/GaoMjun/tcpip => ../tcpip

replace github.com/GaoMjun/goutils => ../goutils

require (
	github.com/GaoMjun/goutils v0.0.0-20230216093315-00bfdc8f95e5
	github.com/GaoMjun/iptransparent v0.0.0-20230216092923-eaa90fe6ac3f
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/golang/snappy v0.0.4
	github.com/gorilla/websocket v1.5.0
	golang.org/x/crypto v0.6.0
)

require (
	github.com/GaoMjun/tcpip v0.0.0-20230216092336-caa5277ffa13 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)
