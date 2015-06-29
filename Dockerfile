FROM golang:1.4-cross

ENV CGO_ENABLED 0

RUN go get gopkg.in/yaml.v2

COPY *.go /go/src/github.com/juliengk/dotfiles/
WORKDIR /go/src/github.com/juliengk/dotfiles

RUN GOARCH=amd64       go build -v -ldflags -d -o /go/bin/dotfiles-amd64
RUN GOARCH=386         go build -v -ldflags -d -o /go/bin/dotfiles-i386
