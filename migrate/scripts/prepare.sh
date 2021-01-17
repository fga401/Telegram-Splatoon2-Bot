#!/bin/bash
# arg 1: path
cd "$(dirname "$0")" || exit 1
cd ..
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
# /bin/bash ./build.sh
echo "Prepare db..."
/bin/bash ./migrate.sh "$data"
echo "Prepare config..."
cp ../config/prod.json "$config"