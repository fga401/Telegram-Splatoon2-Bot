package callbackquery

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"telegram-splatoon2-bot/telegram/bot"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
)

type callbackQueryAdapter struct {
	bot bot.Bot
}

func (a *callbackQueryAdapter) ID() string {
	return "callback"
}

func (a *callbackQueryAdapter) ArgNum() int {
	return 1
}

func (a *callbackQueryAdapter) Adapt(fn adapter.AdaptedFunc, argManager adapter.Manager) adapter.AdaptedFunc {
	argManager.Add(a)
	return func(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
		text := ""
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			err := a.bot.AnswerCallbackQuery(callback.ID)
			if err != nil {
				return errors.Wrap(err, "can't answer callback query")
			}
			text = callbackQueryUtil.GetText(callback.Data)
		}
		return fn(update, argManager, append(args, text)...)
	}
}

// New return a CallbackQuery Adapter which auto answers CallbackQuery and appends text of 'data' to arguments of AdaptedFunc.
func New(bot bot.Bot) adapter.Adapter {
	return &callbackQueryAdapter{
		bot: bot,
	}
}
