#!/bin/bash
addr=$(dig +short myip.opendns.com @resolver1.opendns.com)
cat > src/data/data.json <<EOF
{
	"url": "$addr"
}
EOF
