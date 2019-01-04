#!/bin/sh

rm -rf main/ladder
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o main/ladder main/ladder.go
chmod +x main/ladder

docker image rm mannixgao/ladder
docker build -t mannixgao/ladder .