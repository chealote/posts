#!/bin/bash

curl -s -d '{"name": "che", "secret": "pass"}' localhost:8080/signup
echo "Signed up"

token=$(curl -s -d '{"name": "che", "secret": "pass"}' localhost:8080/signin | jq .token | tr -d \")
echo "Got token: $token"

curl -s -H "Authorization: $token" localhost:8080/
