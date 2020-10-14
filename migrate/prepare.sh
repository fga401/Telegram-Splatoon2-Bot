#!/bin/bash
# arg 1: path
cd "$(dirname "$0")" || exit 1
if [ ! "$1" ]
then
    path=~/bots/splatoon2_bot
else
    path=${1%/}
fi
echo "Path: $path"
data="$path/data"
config="$path/config"
mkdir -p "$data"
mkdir -p "$config"
/bin/bash ./build.sh
/bin/bash ./migrate.sh "$data"
cp ../config/prod.json "$config"