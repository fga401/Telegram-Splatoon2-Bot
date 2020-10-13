#!/bin/bash
# arg 1: version

#git checkout $1
docker build -t splatoon2_bot:"$1" .