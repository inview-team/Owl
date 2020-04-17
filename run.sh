#!/bin/bash

cp .env collector/

kubectl create -f rabbitmq.yaml
kubectl create -f clickhouse.yaml
sleep(100)
kubectl create -f collector/collector.yaml
kubectl create -f worker/worker.yaml
sleep(100)
kubectl create -f metric-analyze/analyzer.yaml
