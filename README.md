# [WIP] [Refactoring] Splatoon2 Telgram Bot
A telegram bot copying all functions in Nintendo app but running on telegram.

The Bot in Version [0.2.0](https://github.com/fga401/Telegram-Splatoon2-Bot/releases/tag/0.2.0) has been deployed, you can find it by [@Splatoon2HelperBot](https://t.me/Splatoon2HelperBot) in Telegram.

## Feature

+ [ ] Salmon query
  + [x] Schedules
  + [ ] Results
+ [x] Stage query
+ [ ] Battle query
  + [ ] Manually
  + [ ] Automatically (whitelist only)
+ [ ] Wiki integration

## Deploy

To build your own bot, clone the codes and generate a new sqlite database file:

```shell script
./migrate/prepare.sh [some_path]
```
The database file and config will be saved in `some_path`. 

Or use your preferred way to execute sqls in `./migrate/sqls/*.up.sql`, and create folders with following structure before running:
```
some_path:
├── config
│   └── prod.json
└── data
    └── sqlite3.db

```

Then build docker image and run:
```shell script
./build.sh [version]
./run.sh [version] [some_path]
```
An example run script:
```shell script
#!/bin/bash
# arg 1: version
# arg 2: path
if [ ! "$1" ]
then
  echo "Need version number!"
  exit 1
fi
if [ ! "$2" ]
then
    path=~/bots/splatoon2_bot
else
    path=${2%/}
fi
echo "Path: "$path
docker stop splatoon2_bot >/dev/null 2>&1
docker rm splatoon2_bot >/dev/null 2>&1
docker run -d -v "$path"/data:/splatoon2_bot/data -v "$path"/config:/splatoon2_bot/config --network host -e socks5_proxy="socks5://127.0.0.1:1080" -e CONFIG=prod -e TOKEN=<token> -e ADMIN=<user_id> -e STORE_CHANNEL=<channel_id> --name splatoon2_bot splatoon2_bot:"$1"
```
You need to change the network and proxy settings, fill the environment values `<token>` and `<channel_id>`, but `-e ADMIN=<user_id>` could be omitted.

If you use proxy, read following docs about how to use proxy in docker:
- https://docs.docker.com/network/proxy/
- https://stackoverflow.com/questions/24319662/from-inside-of-a-docker-container-how-do-i-connect-to-the-localhost-of-the-mach

## Config
You can change the `./config` according to your environment.

- **store channel**: A telegram channel to save some cached images. 
- **others**: ~~I'm too lazy to write README~~ Please read the codes :) 
