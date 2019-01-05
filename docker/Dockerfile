FROM golang:alpine AS builder

RUN go get github.com/GaoMjun/ladder && \
    cd go/src/github.com/GaoMjun/ladder/main && \
    go get -u && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o /ladder

FROM scratch

COPY --from=builder /ladder /bin/ladder

ENTRYPOINT ["ladder"]