package messageutil

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
)

const (
	backTextKey = "« Back to Settings"
)

func AppendBackButton(markup botApi.InlineKeyboardMarkup, target string, printer *message.Printer) botApi.InlineKeyboardMarkup {
	list := markup.InlineKeyboard
	list = append(list,
		botApi.NewInlineKeyboardRow(
			botApi.NewInlineKeyboardButtonData(
				printer.Sprintf(backTextKey),
				callbackQueryUtil.SetPrefix(target, ""),
			),
		),
	)
	return botApi.InlineKeyboardMarkup{
		InlineKeyboard: list,
	}
}
