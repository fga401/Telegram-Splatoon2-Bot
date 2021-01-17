package adapter

import botApi "github.com/go-telegram-bot-api/telegram-bot-api"

type AdaptedFunc func(update botApi.Update, argManager Manager, args ...interface{}) error

type Adapter interface {
	ID() string
	Adapt(fn AdaptedFunc, argManager Manager) AdaptedFunc
	ArgNum() int
}
