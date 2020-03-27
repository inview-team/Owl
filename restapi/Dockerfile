FROM python:3.6

WORKDIR /home/Owl/restapi

COPY . .

RUN pip3 install -r requirements.txt
RUN pip3 install gunicorn

EXPOSE 1337
ENTRYPOINT ["gunicorn", "-b","0.0.0.0:1337", "server"]
