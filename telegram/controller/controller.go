package controller

import botApi "github.com/go-telegram-bot-api/telegram-bot-api"

type Controller interface {
	Handler(update botApi.Update) error
}

