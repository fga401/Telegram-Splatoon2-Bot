#!/bin/bash
cd "$(dirname "$0")" || exit 1
docker build -t splatoon2_bot_migrate -f ./Dockerfile ..
