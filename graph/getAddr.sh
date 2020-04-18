#!/bin/bash
addr=$(kubectl get svc -n owl -o wide | grep restapi | awk '{print $3}')
echo "RESTAPI_SVC_SERVICE_HOST="$addr > .env
