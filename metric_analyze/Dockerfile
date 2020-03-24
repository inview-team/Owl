FROM		ubuntu:16.04
RUN		apt-get update -y && apt-get install -y python3-pip python3-dev
COPY		get_metrics.py /home
COPY		requirements.txt /home
WORKDIR         /home
RUN             apt-get -y update && pip install --upgrade pip && pip3 install -r requirements.txt
ENTRYPOINT	["python3", "get_metrics.py"]
