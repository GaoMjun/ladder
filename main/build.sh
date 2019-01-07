#!/bin/sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o ladder

docker image rm ladder
docker build -t ladder .