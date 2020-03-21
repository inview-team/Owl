FROM		nitincypher/docker-ubuntu-python-pip
FROM   		golang:1.14
ENV     	GO111MODULE on
ENV		GOPATH /go/src
COPY		server.py /home
COPY		requirements.txt /home
WORKDIR         /home
RUN             apt-get -y update && pip install -r requirements.txt
ENTRYPOINT	["python", "server.py"]
