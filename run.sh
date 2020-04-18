#!/bin/bash

kubectl create -f rabbitmq.yaml
kubectl create -f clickhouse.yaml
sleep 30
kubectl create -f collector/collector.yaml
kubectl create -f worker/worker.yaml
sleep 30
kubectl create -f restapi/restapi.yaml
chmod +x graph/getAddr.sh
./graph/getAddr.sh
kubectl create -f graph/front.yaml
kubectl create -f metric-analyze/analyzer.yaml
