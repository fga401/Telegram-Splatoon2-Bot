package message

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
)

const (
	textKeyInternalError = "Internal error. Please retry."
)

// InternalError returns a message showing "internal error".
// It edits the resp message.
func InternalError(printer *message.Printer, resp *botApi.Message) botApi.Chattable {
	text := printer.Sprintf(textKeyInternalError)
	msg := NewByMsg(resp, text, nil, true)
	return msg
}
