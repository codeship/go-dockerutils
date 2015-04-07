FROM golang:1.4.2
MAINTAINER peter.edge@gmail.com

RUN mkdir -p /go/src/github.com/peter-edge/go-dockerutils
ADD . /go/src/github.com/peter-edge/go-dockerutils/
WORKDIR /go/src/github.com/peter-edge/go-dockerutils
