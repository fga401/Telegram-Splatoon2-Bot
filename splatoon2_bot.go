package main

import (
	_ "github.com/mattn/go-sqlite3"
	"telegram-splatoon2-bot/app"
)

func main() {
	app.TelegramApp()
}
