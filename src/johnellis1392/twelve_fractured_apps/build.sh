#!/bin/bash

# Version 1
#
# # Build app binary
# GOOS=linux go build -o app .
#
# # Create docker image
# docker build -t app:v1 .
#
# # Run app
# docker run --rm \
# --volume ./config.json:/etc/config.json \
# --volume ./data:/var/lib/data \
# app:v1

# Version 2
GOOS=linux go build -o app .
docker build -t app:v2 .
docker run --rm \
--env "APP_DATADIR=/var/lib/data" \
--env "APP_HOST=203.0.113.10" \
--env "APP_PORT=3306" \
--env "APP_USERNAME=user" \
--env "APP_PASSWORD=password" \
--env "APP_DATABASE=test" \
app:v2
