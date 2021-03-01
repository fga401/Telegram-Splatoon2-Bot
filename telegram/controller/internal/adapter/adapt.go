package adapter

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/telegram/router"
)

func Apply(fn AdaptedFunc, adapters ...Adapter) router.Handler {
	argManager := NewManager()
	for i := len(adapters) - 1; i >= 0; i-- {
		adapter := adapters[i]
		fn = adapter.Adapt(fn, argManager)
	}
	return func(update botApi.Update) error {
		return fn(update, argManager)
	}
}
