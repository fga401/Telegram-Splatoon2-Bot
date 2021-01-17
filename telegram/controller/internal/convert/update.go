package convert

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/telegram/router"
)

func CallbackQueryToCommand(handler router.Handler) router.Handler {
	return func(update botApi.Update) error {
		update.Message = update.CallbackQuery.Message
		update.Message.From = update.CallbackQuery.From
		update.CallbackQuery = nil
		return handler(update)
	}
}
