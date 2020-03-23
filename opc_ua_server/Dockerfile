FROM		nitincypher/docker-ubuntu-python-pip
COPY		server.py /home
WORKDIR         /home
RUN             apt-get -y update && pip3 install opcua
ENTRYPOINT	["python3", "server.py"]
