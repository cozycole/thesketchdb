#!/bin/bash

mkdir -p ./certbot/conf/live/localhost/

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certbot/conf/live/localhost/privkey.pem \
  -out ./certbot/conf/live/localhost/fullchain.pem \
  -subj "/CN=localhost"
