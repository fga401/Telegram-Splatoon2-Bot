package battle

import (
	"sync"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/nintendo"
	battlePoller "telegram-splatoon2-bot/service/poller/battle"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/bot"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	statusAdapter "telegram-splatoon2-bot/telegram/controller/internal/adapter/status"
	"telegram-splatoon2-bot/telegram/router"
)

// Battle groups all handler about battle result.
type Battle interface {
	BattlePolling(update botApi.Update) error
	BattleAll(update botApi.Update) error
	BattleLast(update botApi.Update) error
	BattleSummary(update botApi.Update) error
}

// UserID is the ID of user
type UserID = userSvc.ID

type battleCtrl struct {
	bot          bot.Bot
	battlePoller battlePoller.Service
	nintendoSvc  nintendo.Service
	userSvc      userSvc.Service
	languageSvc  language.Service

	statusAdapter adapter.Adapter

	battlePollingHandler router.Handler
	battleAllHandler     router.Handler
	battleLastHandler    router.Handler
	battleSummaryHandler router.Handler

	maxResultsPerMessage int
	minLastResults       int

	pollingChats     map[UserID]int64
	pollingMutex     sync.RWMutex
	pollingMaxWorker int32
}

// New returns a Battle object.
func New(bot bot.Bot,
	battlePoller battlePoller.Service,
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

		battlePoller:     battlePoller,
		pollingChats:     make(map[UserID]int64),
		pollingMaxWorker: config.PollingMaxWorker,
	}
	ctrl.battlePollingHandler = adapter.Apply(ctrl.battlePolling, ctrl.statusAdapter)
	ctrl.battleAllHandler = adapter.Apply(ctrl.battleAll, ctrl.statusAdapter)
	ctrl.battleLastHandler = adapter.Apply(ctrl.battleLast, ctrl.statusAdapter)
	ctrl.battleSummaryHandler = adapter.Apply(ctrl.battleSummary, ctrl.statusAdapter)
	go ctrl.pollingRoutine()
	return ctrl
}

func (ctrl *battleCtrl) BattlePolling(update botApi.Update) error {
	return ctrl.battlePollingHandler(update)
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
