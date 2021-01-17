package setting

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/service/language"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/bot"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	callbackQueryAdapter "telegram-splatoon2-bot/telegram/controller/internal/adapter/callbackquery"
	statusAdapter "telegram-splatoon2-bot/telegram/controller/internal/adapter/status"
	"telegram-splatoon2-bot/telegram/router"
)

const (
	KeyboardPrefixReturnToSetting = "<ret_settings>"

	KeyboardPrefixAccountSetting  = "<set_acct>"
	KeyboardPrefixAccountAddition = "<add_acct>"
	KeyboardPrefixAccountDeletion = "<del_acct>"
	KeyboardPrefixAccountSwitch   = "<sw_acct>"
	KeyboardPrefixAccountManager  = "<mgr_acct>"

	KeyboardPrefixLanguageSettings  = "<set_lang>"
	KeyboardPrefixLanguageSelection = "<sel_lang>"

	KeyboardPrefixTimezoneSettings  = "<set_tz>"
	KeyboardPrefixTimezoneSelection = "<sel_tz>"
)

type Setting interface {
	Start(update botApi.Update) error
	Setting(update botApi.Update) error

	LanguageSetting(update botApi.Update) error
	LanguageSelection(update botApi.Update) error

	TimezoneSetting(update botApi.Update) error
	TimezoneSelection(update botApi.Update) error

	AccountSetting(update botApi.Update) error
	AccountManager(update botApi.Update) error
	AccountDeletion(update botApi.Update) error
	AccountAddition(update botApi.Update) error
	AccountSwitch(update botApi.Update) error
}

type settingsCtrl struct {
	bot         bot.Bot
	userSvc     userSvc.Service
	languageSvc language.Service

	callbackQueryAdapter adapter.Adapter
	statusAdapter        adapter.Adapter

	settingHandler router.Handler

	languageSettingHandler   router.Handler
	languageSelectionHandler router.Handler

	timezoneSettingHandler   router.Handler
	timezoneSelectionHandler router.Handler

	accountSettingHandler  router.Handler
	accountManagerHandler  router.Handler
	accountDeletionHandler router.Handler
	accountAdditionHandler router.Handler
	accountSwitchHandler router.Handler
}

func New(bot bot.Bot,
	userSvc userSvc.Service,
	languageSvc language.Service,
) Setting {
	ctrl := &settingsCtrl{
		bot:                  bot,
		userSvc:              userSvc,
		languageSvc:          languageSvc,
		callbackQueryAdapter: callbackQueryAdapter.New(bot),
		statusAdapter:        statusAdapter.New(userSvc),
	}
	ctrl.settingHandler = adapter.Apply(ctrl.setting, ctrl.statusAdapter)

	ctrl.languageSettingHandler = adapter.Apply(ctrl.languageSetting, ctrl.callbackQueryAdapter, ctrl.statusAdapter)
	ctrl.languageSelectionHandler = adapter.Apply(ctrl.languageSelection, ctrl.callbackQueryAdapter, ctrl.statusAdapter)

	ctrl.timezoneSettingHandler = adapter.Apply(ctrl.timezoneSetting, ctrl.callbackQueryAdapter, ctrl.statusAdapter)
	ctrl.timezoneSelectionHandler = adapter.Apply(ctrl.timezoneSelection, ctrl.callbackQueryAdapter, ctrl.statusAdapter)

	ctrl.accountSettingHandler = adapter.Apply(ctrl.accountSetting, ctrl.callbackQueryAdapter, ctrl.statusAdapter)
	ctrl.accountManagerHandler = adapter.Apply(ctrl.accountManager, ctrl.callbackQueryAdapter, ctrl.statusAdapter)
	ctrl.accountDeletionHandler = adapter.Apply(ctrl.accountDeletion, ctrl.callbackQueryAdapter, ctrl.statusAdapter)
	ctrl.accountAdditionHandler = adapter.Apply(ctrl.accountAddition, ctrl.callbackQueryAdapter, ctrl.statusAdapter)
	ctrl.accountSwitchHandler = adapter.Apply(ctrl.accountSwitch, ctrl.callbackQueryAdapter, ctrl.statusAdapter)
	return ctrl
}

func (ctrl *settingsCtrl) Setting(update botApi.Update) error {
	return ctrl.settingHandler(update)
}

func (ctrl *settingsCtrl) TimezoneSetting(update botApi.Update) error {
	return ctrl.timezoneSettingHandler(update)
}

func (ctrl *settingsCtrl) TimezoneSelection(update botApi.Update) error {
	return ctrl.timezoneSelectionHandler(update)
}

func (ctrl *settingsCtrl) LanguageSetting(update botApi.Update) error {
	return ctrl.languageSettingHandler(update)
}

func (ctrl *settingsCtrl) LanguageSelection(update botApi.Update) error {
	return ctrl.languageSelectionHandler(update)
}

func (ctrl *settingsCtrl) AccountSetting(update botApi.Update) error {
	return ctrl.accountSettingHandler(update)
}

func (ctrl *settingsCtrl) AccountManager(update botApi.Update) error {
	return ctrl.accountManagerHandler(update)
}

func (ctrl *settingsCtrl) AccountSwitch(update botApi.Update) error {
	return ctrl.accountSwitchHandler(update)
}

func (ctrl *settingsCtrl) AccountAddition(update botApi.Update) error {
	return ctrl.accountAdditionHandler(update)
}

func (ctrl *settingsCtrl) AccountDeletion(update botApi.Update) error {
	return ctrl.accountDeletionHandler(update)
}