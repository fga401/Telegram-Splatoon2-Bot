package main

//go:generate gotext -srclang=en update -out=locales/catalog.go -lang=en,zh,ja

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"telegram-splatoon2-bot/botutil"
	"telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/nintendo"
	"telegram-splatoon2-bot/service"
)

func InitViper() {
	viper.SetConfigName(os.Getenv("CONFIG"))
	viper.SetConfigType("json")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(errors.Errorf("main: config file: %s \n", err))
	}

	err = viper.BindEnv("token")
	if err != nil {
		panic(errors.Wrap(err, "can't bind token env"))
	}
	err = viper.BindEnv("admin")
	if err != nil {
		panic(errors.Wrap(err, "can't bind admin env"))
	}
	err = viper.BindEnv("store_channel")
	if err != nil {
		panic(errors.Wrap(err, "can't bind store_channel env"))
	}

	if pflag.NArg() == 1 {
		viper.Set("token", pflag.Arg(0))
	}
}

func main() {
	InitViper()
	logger.InitLogger()
	nintendo.InitClient()

	botConfig := botutil.BotConfig{
		UserProxy: viper.GetBool("bot.useProxy"),
		ProxyUrl:  viper.GetString("bot.proxyUrl"),
		Token:     viper.GetString("token"),
		Debug:     viper.GetBool("bot.debug"),
	}
	worker := viper.GetInt("bot.worker")

	myBot := botutil.NewBot(botConfig)
	service.InitService(myBot)

	router := botutil.NewUpdateRouter()
	router.AddCommandHandler("start", service.Start, "Start Command")
	router.AddCommandHandler("settings", service.Settings, "Settings Command")
	router.AddCommandHandler("salmon_schedules", service.QuerySalmonSchedules, "SalmonSchedules Command")
	router.AddCommandHandler("stages", service.QueryStageSchedules, "Stages Command")
	router.AddCallbackQueryHandler(service.AccountSettingsKeyboardPrefix, service.AddAccount, "Settings Callback")
	router.AddCallbackQueryHandler(service.LanguageSettingsKeyboardPrefix, service.SetLanguage, "Settings Callback")
	router.AddCallbackQueryHandler(service.TimezoneSettingsKeyboardPrefix, service.SetTimezone, "Settings Callback")
	router.AddCallbackQueryHandler(service.LanguageSelectionKeyboardPrefix, service.SelectLanguage, "Select Language Callback")
	router.AddCallbackQueryHandler(service.TimezoneSelectionKeyboardPrefix, service.SelectTimezone, "Select Language Callback")
	router.SetTextHandler(service.InputRedirectLink, "Input Redirect Link")

	updateConfig := botapi.UpdateConfig{Offset: 0, Timeout: 60}

	botutil.RunBotInPullMode(myBot, router, updateConfig, worker)
}
