FROM golang

ADD . /go/src/github.com/Xe/gokku

RUN cd /go/src/github.com/Xe/gokku && go get ./...

EXPOSE 80
ENV PORT 80

RUN apt-get update && apt-get install openssh-client -y

RUN mkdir -p /gokku/repo

CMD gokku
