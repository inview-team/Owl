FROM		nitincypher/docker-ubuntu-python-pip
COPY		get_metrics.py /home
WORKDIR         /home
RUN             apt-get -y update && pip install clickhouse_driver 
ENTRYPOINT	["python", "get_metrics.py"]
