package router

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/common/enum"
)

type Option enum.Enum
type optionEnum struct {
	CaseSensitive Option
	AsDefault     Option
}

var OptionEnum = enum.Assign(&optionEnum{}).(*optionEnum)

type Handler func(message botApi.Update) error

type Router interface {
	RegisterCommand(command string, handler Handler, options ...Option)
	RegisterCallbackQuery(prefix string, handler Handler)
	RegisterText(handler Handler)

	Run()
}
