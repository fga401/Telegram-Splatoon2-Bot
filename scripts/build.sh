#!/bin/bash
# arg 1: docker tag
# arg 2: git branch or tag
if [ ! "$1" ]
then
    version="latest"
    git checkout master
else
    version=$1
    git_version=$1
    if [ "$2" ]
    then
        git_version=$2
    fi
    git checkout "$git_version"
fi
cd "$(dirname "$0")" || exit 1
cd ..
docker build -t splatoon2_bot:"$version" .