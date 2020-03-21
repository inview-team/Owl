FROM		nitincypher/docker-ubuntu-python-pip
COPY		server.py /home
COPY		requirements.txt /home
WORKDIR         /home
RUN             apt-get -y update && pip install -r requirements.txt
ENTRYPOINT	["python", "server.py"]
