FROM yikaus/alpine-base

RUN apk update && apk add go openssh git bash

EXPOSE 80
ENV PORT 80
ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN mkdir -p /gokku/repo

ADD . /go/src/github.com/Xe/gokku

RUN cd /go/src/github.com/Xe/gokku && go get ./...

CMD gokku
