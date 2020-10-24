#!/bin/bash
# arg 1: version
if [ ! "$1" ]
then
    version="latest"
else
    version=$1
    git checkout "$version"
fi
cd "$(dirname "$0")" || exit 1
docker build -t splatoon2_bot:"$version" .