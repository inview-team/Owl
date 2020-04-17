FROM    	golang:1.13-alpine3.10
ENV     	GO111MODULE on
COPY		collector.go /home
WORKDIR		/home
ENTRYPOINT	["go", "run", "collector.go"]
