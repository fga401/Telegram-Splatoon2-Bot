#!/bin/bash
cd "$(dirname "$0")" || exit 1
cd ..
docker build -t splatoon2_bot_migrate -f ./Dockerfile ..
