#!/bin/sh
echo "RESTAPI_SVC_SERVICE_HOST=" > /tmp/nginx/frontend/.env
addr=$(echo $RESTAPI_SVC_SERVICE_HOST)
sed '$ s/$/'$addr'/' /tmp/nginx/frontend/.env > /tmp/nginx/frontend/.env.new && mv /tmp/nginx/frontend/.env.new /tmp/nginx/frontend/.env
nginx -g "daemon off;"
