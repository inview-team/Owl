#!/bin/bash

docker run -d --name ch-db  --ulimit nofile=262144:262144 -p 9000:9000 yandex/clickhouse-server
