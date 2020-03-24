FROM python:3.6-alpine

WORKDIR /home/Owl/restapi

COPY . .

RUN python3 -m venv venv
RUN venv/bin/pip3 install -r requirements.txt
RUN venv/bin/pip3 install gunicorn


EXPOSE 5000
CMD ["python3", "server.py"]
