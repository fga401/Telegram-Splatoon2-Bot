package message

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// NewByMsg returns a telegram message given text, markup and msg.
// ChatID would be set as msg.
func NewByMsg(msg *botApi.Message, text string, markup *botApi.InlineKeyboardMarkup, edit bool) botApi.Chattable {
	if edit {
		ret := botApi.NewEditMessageText(msg.Chat.ID, msg.MessageID, text)
		ret.ParseMode = "Markdown"
		if markup != nil {
			ret.ReplyMarkup = markup
		}
		return ret
	}
	ret := botApi.NewMessage(msg.Chat.ID, text)
	ret.ParseMode = "Markdown"
	if markup != nil {
		ret.ReplyMarkup = &markup
	}
	return ret
}

// NewByUpdate returns a telegram message given text, markup and update.
// If it's a CallbackQuery update, the new message will edit origin message. Otherwise, it will generate a new message.
func NewByUpdate(update botApi.Update, text string, markup *botApi.InlineKeyboardMarkup) botApi.Chattable {
	if update.CallbackQuery != nil {
		return NewByMsg(update.CallbackQuery.Message, text, markup, true)
	}
	return NewByMsg(update.Message, text, markup, false)
}
