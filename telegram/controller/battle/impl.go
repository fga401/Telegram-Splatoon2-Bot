package battle

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/nintendo"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/bot"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	statusAdapter "telegram-splatoon2-bot/telegram/controller/internal/adapter/status"
	"telegram-splatoon2-bot/telegram/router"
)

// Battle groups all handler about battle result.
type Battle interface {
	BattleAll(update botApi.Update) error
	BattleLast(update botApi.Update) error
	BattleSummary(update botApi.Update) error
}

type battleCtrl struct {
	bot         bot.Bot
	nintendoSvc nintendo.Service
	userSvc     userSvc.Service
	languageSvc language.Service

	statusAdapter adapter.Adapter

	battleAllHandler     router.Handler
	battleLastHandler    router.Handler
	battleSummaryHandler router.Handler

	maxResultsPerMessage int
	minLastResults       int
}

// New returns a Battle object.
func New(bot bot.Bot,
	nintendoSvc nintendo.Service,
	userSvc userSvc.Service,
	languageSvc language.Service,
	config Config,
) Battle {
	ctrl := &battleCtrl{
		bot:           bot,
		nintendoSvc:   nintendoSvc,
		userSvc:       userSvc,
		languageSvc:   languageSvc,
		statusAdapter: statusAdapter.New(userSvc),

		maxResultsPerMessage: config.MaxResultsPerMessage,
		minLastResults:       config.MinLastResults,
	}
	ctrl.battleAllHandler = adapter.Apply(ctrl.battleAll, ctrl.statusAdapter)
	ctrl.battleLastHandler = adapter.Apply(ctrl.battleLast, ctrl.statusAdapter)
	ctrl.battleSummaryHandler = adapter.Apply(ctrl.battleSummary, ctrl.statusAdapter)
	return ctrl
}

func (ctrl *battleCtrl) BattleAll(update botApi.Update) error {
	return ctrl.battleAllHandler(update)
}

func (ctrl *battleCtrl) BattleLast(update botApi.Update) error {
	return ctrl.battleLastHandler(update)
}

func (ctrl *battleCtrl) BattleSummary(update botApi.Update) error {
	return ctrl.battleSummaryHandler(update)
}
