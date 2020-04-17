FROM		nitincypher/docker-ubuntu-python-pip
COPY		server.py /home
WORKDIR         /home
RUN             apt-get -y update && pip install opcua
ENTRYPOINT	["python", "server.py"]
