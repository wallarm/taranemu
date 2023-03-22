#!/usr/bin/bash -eu

docker run \
    --name tarantool \
    -d -p 3301:3301 \
    -v /var/lib/tarantool:/var/lib/tarantool \
    -e TARANTOOL_USER_NAME=user \
    -e TARANTOOL_USER_PASSWORD=DSoXbver3p4bbMK6dGhUfo \
    --restart=unless-stopped \
    tarantool/tarantool:latest
