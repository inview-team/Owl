FROM    	golang:1.13-alpine3.10
ENV     	GO111MODULE on
RUN     	apk update && apk upgrade
COPY		collector.go $GOPATH
ENTRYPOINT	["go", "run", "collector.go"]
