package app

import (
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
	proxyClient "telegram-splatoon2-bot/common/proxyclient"
	"telegram-splatoon2-bot/driver/cache/fastcache"
	"telegram-splatoon2-bot/driver/cache/gocache"
	"telegram-splatoon2-bot/driver/database"
	imgDownloader "telegram-splatoon2-bot/service/image/downloader"
	tgImageUploader "telegram-splatoon2-bot/service/image/uploader/telegram"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/repository"
	"telegram-splatoon2-bot/service/repository/salmon"
	"telegram-splatoon2-bot/service/repository/stage"
	"telegram-splatoon2-bot/service/timezone"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/bot"
	"telegram-splatoon2-bot/telegram/controller/battle"
	repositoryCtrl "telegram-splatoon2-bot/telegram/controller/repository"
	"telegram-splatoon2-bot/telegram/router"
)

func init() {
	// load from file
	viper.SetConfigName(os.Getenv("CONFIG"))
	viper.SetConfigType("json")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("can't read config file", zap.Error(err))
	}
	// load from environment variables
	err = viper.BindEnv("token")
	if err != nil {
		log.Panic("can't bind token env", zap.Error(err))
	}
	err = viper.BindEnv("admin")
	if err != nil {
		log.Panic("can't bind admin env", zap.Error(err))
	}
	err = viper.BindEnv("store_channel")
	if err != nil {
		log.Panic("can't bind store_channel env", zap.Error(err))
	}
	// load from CLI arguments
	if pflag.NArg() == 1 {
		viper.Set("token", pflag.Arg(0))
	}
	// init logger
	log.InitLogger(logConfig())
}

func logConfig() string {
	return viper.GetString("log.level")
}

func token() string {
	return viper.GetString("token")
}

func botAPIClientConfig() proxyClient.Config {
	return proxyClient.Config{
		EnableProxy: viper.GetBool("bot.client.enableProxy"),
		ProxyURL:    viper.GetString("bot.client.proxyURL"),
		EnableHTTP2: true,
		Timeout:     viper.GetDuration("bot.client.timeout"),
	}
}

func botConfig() bot.Config {
	return bot.Config{
		RetryTimes: viper.GetInt("bot.retryTimes"),
		DefaultCallbackQueryConfig: bot.CallbackQueryConfig{
			Text:      "",
			ShowAlert: false,
			CacheTime: viper.GetInt("bot.callBackQuery.cacheTimeInSecond"),
		},
	}
}

func routerConfig() router.Config {
	mode := router.ModeEnum.Polling
	modeStr := strings.ToLower(viper.GetString("router.mode"))
	if modeStr == "webhook" {
		mode = router.ModeEnum.WebHook
	}
	return router.Config{
		Mode:      mode,
		MaxWorker: viper.GetInt32("router.maxWorker"),
		Polling: router.PollingConfig{
			Timeout: viper.GetInt("router.polling.timeoutInSecond"),
		},
	}
}

func databaseConfig() database.Config {
	return database.Config{
		URL:          viper.GetString("database.url"),
		Driver:       viper.GetString("database.driver"),
		MaxIdleConns: viper.GetInt("database.maxIdleConns"),
		MaxOpenConns: viper.GetInt("database.maxOpenConns"),
	}
}

func fastcacheConfig() fastcache.Config {
	return fastcache.Config{
		MaxBytes: viper.GetInt("fastcache.maxBytes"),
	}
}

func proofKeyCacheConfig() gocache.Config {
	return gocache.Config{
		Expiration: viper.GetDuration("gocache.proofKey.expiration"),
		CleanUp:    viper.GetDuration("gocache.proofKey.cleanUp"),
	}
}

func nintendoConfig() nintendo.Config {
	return nintendo.Config{
		Timeout:    viper.GetDuration("nintendo.client.timeout"),
		RetryTimes: viper.GetInt("nintendo.retryTimes"),
	}
}

func userSvcConfig() userSvc.Config {
	adminsConfig := viper.GetStringSlice("admin")
	adminsID := make([]userSvc.ID, 0, len(adminsConfig))
	for _, s := range adminsConfig {
		id, err := strconv.Atoi(s)
		if err != nil {
			log.Panic("can't load admin", zap.Error(err))
		}
		adminsID = append(adminsID, userSvc.ID(id))
	}
	return userSvc.Config{
		DefaultPermission: userSvc.DefaultPermission{
			Admins:       adminsID,
			MaxAccount:   viper.GetInt32("user.permission.maxAccount"),
			AllowPolling: viper.GetBool("user.permission.allowPolling"),
			Timezone:     timezone.ByMinute(viper.GetInt("user.permission.timezone")),
			Language:     language.ByIETF(viper.GetString("user.permission.language")),
			IsBlock:      false,
		},
		AccountsCacheExpiration: viper.GetDuration("user.accountExpiration"),
		ProofKeyCacheExpiration: viper.GetDuration("user.proofKeyExpiration"),
	}
}

func languageSvcConfig() language.Config {
	return language.Config{
		SupportedLanguages: viper.GetStringSlice("language"),
		LocalePath:         viper.GetString("locale.path"),
	}
}

func tgImgUploaderConfig() tgImageUploader.Config {
	storeChannelID, err := strconv.ParseInt(viper.GetString("store_channel"), 10, 64)
	if err != nil {
		log.Panic("can't parse Image Config: StoreChannelID", zap.Error(err))
	}
	return tgImageUploader.Config{
		StoreChannelID: storeChannelID,
	}
}

func imgDownloaderConfig() imgDownloader.Config {
	return imgDownloader.Config{
		Proxy: proxyClient.Config{
			EnableProxy: false,
			ProxyURL:    "",
			EnableHTTP2: false,
			Timeout:     0,
		},
		RetryTimes: viper.GetInt("image.retryTimes"),
	}
}

func salmonRepositoryConfig() salmon.Config {
	return salmon.Config{
		Dumper: salmon.DumperConfig{
			StageFile:  viper.GetString("repository.salmon.stageFileName"),
			WeaponFile: viper.GetString("repository.salmon.weaponFileName"),
		},
		RandomWeaponPath:  viper.GetString("repository.salmon.randomWeaponImagePath"),
		GrizzcoWeaponPath: viper.GetString("repository.salmon.grizzcoWeaponImagePath"),
	}
}

func stageRepositoryConfig() stage.Config {
	return stage.Config{
		Dumper: stage.DumperConfig{
			StageFile: viper.GetString("repository.stage.stageFileName"),
		},
	}
}

func repositoryManagerConfig() repository.ManagerConfig {
	return repository.ManagerConfig{
		Delay: viper.GetDuration("repository.delay"),
	}
}

func repositoryControllerConfig() repositoryCtrl.Config {
	return repositoryCtrl.Config{
		Limit: viper.GetInt("controller.limit"),
	}
}

func battleControllerConfig() battle.Config {
	return battle.Config{
		MaxResultsPerMessage: viper.GetInt("controller.maxBattleResultsPerMessage"),
		MinLastResults:       viper.GetInt("controller.minLastBattleResults"),
	}
}
