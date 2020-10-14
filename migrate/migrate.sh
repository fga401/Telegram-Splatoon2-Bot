#!/bin/bash
# arg 1: db file path
cd "$(dirname "$0")" || exit 1
if [ ! "$1" ]
then
    data_path=~/sqlite/splatoon2_bot
else
    data_path=${1%/}
fi
docker stop splatoon2_bot_migrate >/dev/null 2>&1
docker rm splatoon2_bot_migrate >/dev/null 2>&1
docker run -v "$data_path":/splatoon2_bot/data -v "$(pwd)"/sqls:/splatoon2_bot/migrate/sqls -e CONFIG=prod --name splatoon2_bot_migrate splatoon2_bot_migrate
