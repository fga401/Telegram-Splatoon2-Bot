package repository

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/repository/salmon"
	"telegram-splatoon2-bot/service/repository/stage"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/bot"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	callbackQueryAdapter "telegram-splatoon2-bot/telegram/controller/internal/adapter/callbackquery"
	statusAdapter "telegram-splatoon2-bot/telegram/controller/internal/adapter/status"
	"telegram-splatoon2-bot/telegram/router"
)

// Repository groups all handler about schedules.
type Repository interface {
	Salmon(update botApi.Update) error
	Stage(update botApi.Update) error
}

type repositoryCtrl struct {
	bot         bot.Bot
	userSvc     userSvc.Service
	languageSvc language.Service

	salmonRepo salmon.Repository
	stageRepo  stage.Repository

	callbackQueryAdapter adapter.Adapter
	statusAdapter        adapter.Adapter

	salmonHandler router.Handler
	stageHandler  router.Handler

	limit int
}

// New returns a Repository object.
func New(bot bot.Bot,
	userSvc userSvc.Service,
	languageSvc language.Service,
	salmonRepo salmon.Repository,
	stageRepo stage.Repository,
	config Config,
) Repository {
	ctrl := &repositoryCtrl{
		bot:         bot,
		userSvc:     userSvc,
		languageSvc: languageSvc,

		callbackQueryAdapter: callbackQueryAdapter.New(bot),
		statusAdapter:        statusAdapter.New(userSvc),

		salmonRepo: salmonRepo,
		stageRepo:  stageRepo,

		limit: config.Limit,
	}
	ctrl.salmonHandler = adapter.Apply(ctrl.salmon, ctrl.callbackQueryAdapter, ctrl.statusAdapter)
	ctrl.stageHandler = adapter.Apply(ctrl.stage, ctrl.callbackQueryAdapter, ctrl.statusAdapter)
	return ctrl
}

func (ctrl *repositoryCtrl) Salmon(update botApi.Update) error {
	return ctrl.salmonHandler(update)
}

func (ctrl *repositoryCtrl) Stage(update botApi.Update) error {
	return ctrl.stageHandler(update)
}
