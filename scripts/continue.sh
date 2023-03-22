#!/bin/bash -eu

cd $(dirname $(realpath $0))/..

while true;
do
    echo "Running tarantella-server..."
    env $(cat .env | xargs) go run ./tarantella-server/ || sleep 3s
done
