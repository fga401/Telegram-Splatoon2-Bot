package help

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/service/language"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/bot"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	statusAdapter "telegram-splatoon2-bot/telegram/controller/internal/adapter/status"
	"telegram-splatoon2-bot/telegram/router"
)

// Help groups all handler about help.
type Help interface {
	Help(update botApi.Update) error
	HelpStages(update botApi.Update) error
}

type helpCtrl struct {
	bot         bot.Bot
	userSvc     userSvc.Service
	languageSvc language.Service

	statusAdapter adapter.Adapter

	helpHandler       router.Handler
	helpStagesHandler router.Handler
}

// New returns a Help object.
func New(bot bot.Bot,
	userSvc userSvc.Service,
	languageSvc language.Service,
) Help {
	ctrl := &helpCtrl{
		bot:           bot,
		userSvc:       userSvc,
		languageSvc:   languageSvc,
		statusAdapter: statusAdapter.New(userSvc),
	}
	ctrl.helpHandler = adapter.Apply(ctrl.help, ctrl.statusAdapter)
	ctrl.helpStagesHandler = adapter.Apply(ctrl.helpStages, ctrl.statusAdapter)
	return ctrl
}

func (ctrl *helpCtrl) Help(update botApi.Update) error {
	return ctrl.helpHandler(update)
}

func (ctrl *helpCtrl) HelpStages(update botApi.Update) error {
	return ctrl.helpStagesHandler(update)
}
