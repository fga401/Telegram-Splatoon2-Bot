package app

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
	proxyClient "telegram-splatoon2-bot/common/proxyclient"
	"telegram-splatoon2-bot/driver/cache/fastcache"
	"telegram-splatoon2-bot/driver/cache/gocache"
	"telegram-splatoon2-bot/driver/cache/syncmap"
	"telegram-splatoon2-bot/driver/database"
	imageSvc "telegram-splatoon2-bot/service/image"
	imgDownloader "telegram-splatoon2-bot/service/image/downloader"
	tgImgUploader "telegram-splatoon2-bot/service/image/uploader/telegram"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/repository"
	"telegram-splatoon2-bot/service/repository/salmon"
	"telegram-splatoon2-bot/service/repository/stage"
	userSvc "telegram-splatoon2-bot/service/user"
	userDatabase "telegram-splatoon2-bot/service/user/database"
	"telegram-splatoon2-bot/telegram/bot"
	repositoryCtrl "telegram-splatoon2-bot/telegram/controller/repository"
	"telegram-splatoon2-bot/telegram/controller/setting"
	"telegram-splatoon2-bot/telegram/router"
	"telegram-splatoon2-bot/telegram/controller/help"
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
	accountCache := fastcache.New(fastcacheConfig())
	proofKeyCache := gocache.New(proofKeyCacheConfig())

	nintendoSvc := nintendo.New(nintendoConfig())

	userSvc := userSvc.New(userDatabase, adminCache, statusCache, accountCache, proofKeyCache, nintendoSvc, userSvcConfig())
	languageSvc := language.NewService(languageSvcConfig())

	imgUploader := tgImgUploader.NewUploader(bot, tgImgUploaderConfig())
	imgDownloader := imgDownloader.NewDownloader(imgDownloaderConfig())
	imageSvc := imageSvc.NewService(imgUploader, imgDownloader)
	salmonRepo := salmon.NewRepository(nintendoSvc, userSvc, imageSvc, salmonRepositoryConfig())
	stageRepo := stage.NewRepository(nintendoSvc, userSvc, imageSvc, stageRepositoryConfig())
	repoManager := repository.NewManager(repositoryManagerConfig(), salmonRepo, stageRepo)
	repoManager.Start()

	settingCtrl := setting.New(bot, userSvc, languageSvc)
	router.RegisterCommand("start", settingCtrl.Start)
	router.RegisterCommand("settings", settingCtrl.Setting)
	router.RegisterCallbackQuery(setting.KeyboardPrefixSetting, settingCtrl.Setting)
	router.RegisterCallbackQuery(setting.KeyboardPrefixCancelSetting, settingCtrl.CancelSetting)
	router.RegisterCallbackQuery(setting.KeyboardPrefixLanguageSettings, settingCtrl.LanguageSetting)
	router.RegisterCallbackQuery(setting.KeyboardPrefixLanguageSelection, settingCtrl.LanguageSelection)
	router.RegisterCallbackQuery(setting.KeyboardPrefixTimezoneSettings, settingCtrl.TimezoneSetting)
	router.RegisterCallbackQuery(setting.KeyboardPrefixTimezoneSelection, settingCtrl.TimezoneSelection)
	router.RegisterCallbackQuery(setting.KeyboardPrefixAccountSetting, settingCtrl.AccountSetting)
	router.RegisterCallbackQuery(setting.KeyboardPrefixAccountSwitch, settingCtrl.AccountSwitch)
	router.RegisterCallbackQuery(setting.KeyboardPrefixAccountManager, settingCtrl.AccountManager)
	router.RegisterCallbackQuery(setting.KeyboardPrefixAccountAddition, settingCtrl.AccountAddition)
	router.RegisterCallbackQuery(setting.KeyboardPrefixAccountDeletionConfirm, settingCtrl.AccountDeletionConfirm)
	router.RegisterCallbackQuery(setting.KeyboardPrefixAccountDeletion, settingCtrl.AccountDeletion)
	router.RegisterText(settingCtrl.AccountRedirectLink)

	repoCtrl := repositoryCtrl.New(bot, userSvc, languageSvc, salmonRepo, stageRepo, repositoryControllerConfig())
	router.RegisterCommand("salmon_schedules", repoCtrl.Salmon)
	router.RegisterCommand("stages", repoCtrl.Stage)

	helpCtrl := help.New(bot, userSvc, languageSvc)
	router.RegisterCommand("help", helpCtrl.Help)
	router.RegisterCommand("help_stages", helpCtrl.HelpStages)

	router.Run()
}
