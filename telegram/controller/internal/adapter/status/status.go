package status

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
)

type statusAdapter struct {
	userSvc userSvc.Service
}

func (a *statusAdapter) ID() string {
	return "status"
}

func (a *statusAdapter) ArgNum() int {
	return 1
}

func (a *statusAdapter) Adapt(fn adapter.AdaptedFunc, argManager adapter.Manager) adapter.AdaptedFunc {
	argManager.Add(a)
	return func(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
		var user *botApi.User
		if update.Message != nil && update.Message.From != nil {
			user = update.Message.From
		} else if update.CallbackQuery != nil && update.CallbackQuery.From != nil {
			user = update.CallbackQuery.From
		} else {
			return errors.New("user not found in update or unsupported update type")
		}
		userID := userSvc.ID(user.ID)
		status, err := a.userSvc.GetStatus(userID)
		if err != nil {
			return errors.Wrap(err, "can't fetch status")
		}
		return fn(update, argManager, append(args, status)...)
	}
}

func New(userSvc userSvc.Service) adapter.Adapter {
	return &statusAdapter{
		userSvc: userSvc,
	}
}
