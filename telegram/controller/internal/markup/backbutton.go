package markup

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
)

const (
	textKeyBack = "« Go Back"
)

func AppendBackButton(markup botApi.InlineKeyboardMarkup, target string, printer *message.Printer) botApi.InlineKeyboardMarkup {
	list := markup.InlineKeyboard
	list = append(list,
		botApi.NewInlineKeyboardRow(
			botApi.NewInlineKeyboardButtonData(
				printer.Sprintf(textKeyBack),
				callbackQueryUtil.SetPrefix(target, ""),
			),
		),
	)
	return botApi.InlineKeyboardMarkup{
		InlineKeyboard: list,
	}
}
