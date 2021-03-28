package router

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/common/enum"
)

//Option uses in Router.RegisterCommand to provide more options.
type Option enum.Enum
type optionEnum struct {
	// CaseSensitive indicates that the command is case sensitive.
	CaseSensitive Option
	// AsDefault indicates that use this handler if no command matched.
	AsDefault Option
	// Regexp indicates that the command is a regular expression.
	// All command matching this regular expression would be processed by handler.
	// Earlier registered commands take higher priority.
	Regexp Option
}

// OptionEnum lists all available option.
var OptionEnum = enum.Assign(&optionEnum{}).(*optionEnum)

// Handler is a function that process request from telegram.
type Handler func(message botApi.Update) error

// Router manages a batch of handlers to process different input form telegram.
type Router interface {
	// RegisterCommand adds a handler to process '/command' message.
	RegisterCommand(command string, handler Handler, options ...Option)
	// RegisterCallbackQuery adds a handler to CallbackQuery request.
	// All CallbackQuery with this prefix would be processed by this handler.
	// Prefix is a substring of the 'data' CallbackQuery field ending at the first colon. e.g. "<sw_acct>" is the prefix of "<sw_acct>:xxx".
	RegisterCallbackQuery(prefix string, handler Handler)
	// RegisterText adds a handler to process plain text (not a command) input.
	RegisterText(handler Handler)

	// Run starts the Router.
	Run()
}
