package app

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
	proxyClient "telegram-splatoon2-bot/common/proxyclient"
	"telegram-splatoon2-bot/driver/cache/fastcache"
	"telegram-splatoon2-bot/driver/database"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/bot"
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

func botAPiClientConfig() proxyClient.Config {
	return proxyClient.Config{
		EnableProxy: viper.GetBool("bot.client.enableProxy"),
		ProxyUrl:    "",
		EnableHttp2: false,
		Timeout: viper.GetDuration("bot.client.timeout"),
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

func userSvcConfig() userSvc.Config {
	adminsConfig := viper.GetIntSlice("admin")
	adminsID := make([]userSvc.ID, 0, len(adminsConfig))
	for _, id := range adminsConfig {
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
	}
}

func languageSvcConfig() language.Config {
	return language.Config{
		SupportedLanguages: viper.GetStringSlice("language"),
	}
}