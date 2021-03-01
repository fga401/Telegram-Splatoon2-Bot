package adapter

import botApi "github.com/go-telegram-bot-api/telegram-bot-api"

// AdaptedFunc is a function wrapping Handler to provide more custom arguments.
type AdaptedFunc func(update botApi.Update, argManager Manager, args ...interface{}) error

// Adapter adds more arguments to AdaptedFunc.
type Adapter interface {
	// ID is the unique key of Adapter.
	ID() string
	// Adapt appends new arguments to the args of AdaptedFunc.
	Adapt(fn AdaptedFunc, argManager Manager) AdaptedFunc
	// ArgNum returns the number of new arguments which this Adapter adds.
	ArgNum() int
}
