package message

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewByMsg(msg *botApi.Message, text string, markup *botApi.InlineKeyboardMarkup, edit bool) botApi.Chattable {
	if edit {
		ret := botApi.NewEditMessageText(msg.Chat.ID, msg.MessageID, text)
		ret.ParseMode = "Markdown"
		if markup != nil {
			ret.ReplyMarkup = markup
		}
		return ret
	} else {
		ret := botApi.NewMessage(msg.Chat.ID, text)
		ret.ParseMode = "Markdown"
		if markup != nil {
			ret.ReplyMarkup = &markup
		}
		return ret
	}
}

func NewByUpdate(update botApi.Update, text string, markup *botApi.InlineKeyboardMarkup) botApi.Chattable {
	if update.CallbackQuery != nil {
		return NewByMsg(update.CallbackQuery.Message, text, markup, true)
	} else {
		return NewByMsg(update.Message, text, markup, false)
	}
}
