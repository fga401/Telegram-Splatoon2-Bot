package app

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
	proxyClient "telegram-splatoon2-bot/common/proxyclient"
	"telegram-splatoon2-bot/driver/cache/fastcache"
	"telegram-splatoon2-bot/driver/cache/syncmap"
	"telegram-splatoon2-bot/driver/database"
	"telegram-splatoon2-bot/service/language"
	userSvc "telegram-splatoon2-bot/service/user"
	userDatabase "telegram-splatoon2-bot/service/user/database"
	"telegram-splatoon2-bot/telegram/bot"
	"telegram-splatoon2-bot/telegram/controller/setting"
	"telegram-splatoon2-bot/telegram/router"
)

func TelegramApp() {
	botClient := proxyClient.New(botAPiClientConfig())
	botApi, err := botApi.NewBotAPIWithClient(token(), botClient)
	if err != nil {
		log.Panic("can't init botApi", zap.Error(err))
	}
	bot := bot.New(botApi, botConfig())
	router := router.New(botApi, routerConfig())

	database := database.New(databaseConfig())
	userDatabase := userDatabase.New(database)
	adminCache := syncmap.New()
	statusCache := fastcache.New(fastcacheConfig())

	userSvc := userSvc.New(userDatabase, adminCache, statusCache, userSvcConfig())
	languageSvc := language.NewService(languageSvcConfig())
	settingCtrl := setting.New(bot, userSvc, languageSvc)

	router.RegisterCommand("start", settingCtrl.Start)
	router.RegisterCommand("settings", settingCtrl.Setting)
	router.RegisterCallbackQuery(setting.KeyboardPrefixLanguageSettings, settingCtrl.LanguageSetting)
	router.RegisterCallbackQuery(setting.KeyboardPrefixLanguageSelection, settingCtrl.LanguageSelection)
	router.RegisterCallbackQuery(setting.KeyboardPrefixTimezoneSettings, settingCtrl.TimezoneSetting)
	router.RegisterCallbackQuery(setting.KeyboardPrefixTimezoneSelection, settingCtrl.TimezoneSelection)

	router.Run()
}
