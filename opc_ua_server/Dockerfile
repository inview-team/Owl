FROM		ubuntu:16.04
RUN		apt-get update -y && apt-get install -y python3-pip python3-dev
COPY		server.py /home
WORKDIR         /home
RUN             apt-get -y update && pip3 install opcua
ENTRYPOINT	["python3", "server.py"]
