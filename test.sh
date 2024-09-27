#!/bin/bash

set -e

rm posts.db

./posts -i
./posts &
pid="$!"

username="che"
password="somepassword"

echo "Login information:"
echo "  username: $username"
echo "  password: $password"

curl -s -d '{"name": "che", "password": "somepassword"}' localhost:8080/signup
echo "Signed up"

token="invalid"
echo "Using token '$token' to access /"
curl -s -H "Authorization: $token" localhost:8080/

token=$(curl -s -d '{"name": "che", "password": "somepassword"}' localhost:8080/signin | jq .token | tr -d \")
echo "Logged in. Got token: $token"

echo "Using token '$token' to access /"
curl -s -H "Authorization: $token" localhost:8080/
echo ""

echo "Using token '$token' to access /somepage"
curl -v -s -H "Authorization: $token" localhost:8080/somepage
echo ""

echo "Sleeping 5 seconds and trying again..."
kill $pid
sleep 5

./posts &
pid="$!"
echo "Using token '$token' to access /"
curl -s -H "Authorization: $token" localhost:8080/
echo ""

echo "Sleeping 15 seconds and trying again..."
kill $pid
sleep 15

./posts &
pid="$!"
echo "Using token '$token' to access /"
curl -s -H "Authorization: $token" localhost:8080/
echo ""

kill $pid
