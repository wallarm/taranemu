#!/bin/bash -eu

cd $(dirname $(realpath $0))/..
echo -en "\033]0;⏱🤖 Watching $(pwd)...\a"
#⌚

BINDIR=$(realpath .bin)

while inotifywait -e close_write -r . --exclude '(\.git)|(testdata)|(\.bin)|(console)' ;
do
    echo -en "\033]0;⏯ Sending stop signal!\a"
    killall -TERM tarantella-server || true
    echo -en "\033]0;⏱ Watching $(pwd)...\a"
done
