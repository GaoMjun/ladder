#!/bin/sh

rm -rf ladder
rm -rf ladder_fc.zip

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o ladder

upx -9 -v ladder

zip ladder_fc.zip ladder

