package message

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
)

const (
	textKeyUpdatingToken = "Your token is expired. Fetching new token from Nintendo server..."
)

// UpdatingToken returns a message showing "updating token"
func UpdatingToken(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyUpdatingToken)
	msg := NewByUpdate(update, text, nil)
	return msg
}
