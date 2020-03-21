FROM   		golang:1.13-alpine3.10
ENV     	GO111MODULE on
ENV		GOPATH /go/src
RUN     	apk update && apk upgrade && apk add vim bash
COPY		collector.go /home
WORKDIR		/home
ENTRYPOINT	["go", "run", "collector.go"]
