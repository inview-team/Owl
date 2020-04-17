FROM    	golang:1.13-alpine3.10
ENV     	GO111MODULE on
COPY		analyzer.go /home
COPY		databaseQuery.go /home
COPY		sendNotification.go /home
COPY		.env /home
WORKDIR		/home
ENTRYPOINT	["go", "run", "."]
