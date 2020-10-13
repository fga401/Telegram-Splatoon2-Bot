#!/bin/bash
# arg 1: db file path
bash_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
if [ ! "$1" ]
then
    data_path=~/sqlite/splatoon2_bot/data
else
    data_path=$1
fi
docker stop splatoon2_bot_migrate >/dev/null 2>&1
docker rm splatoon2_bot_migrate >/dev/null 2>&1
docker build -t splatoon2_bot_migrate "$bash_dir"/../
docker run -v "$data_path":/splatoon2/data -e CONFIG=prod --name splatoon2_bot_migrate splatoon2_bot_migrate
