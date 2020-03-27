FROM    	golang:1.13-alpine3.10
ENV     	GO111MODULE on
RUN     	apk update && apk upgrade
COPY		worker.go $GOPATH
ENTRYPOINT	["go", "run", "worker.go"]
