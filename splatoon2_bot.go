package main

//go:generate gotext -srclang=en update -out=locales/catalog.go -lang=en,zh,ja

//import (
//	"os"
//
//	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
//	"github.com/pkg/errors"
//	"github.com/spf13/pflag"
//	"github.com/spf13/viper"
//	"go.uber.org/zap"
//	"telegram-splatoon2-bot/botutil"
//	"telegram-splatoon2-bot/common/log"
//	"telegram-splatoon2-bot/service/nintendo"
//	"telegram-splatoon2-bot/service/repository/salmon"
//	"telegram-splatoon2-bot/service/repository/stage"
//	"telegram-splatoon2-bot/service/todo"
//)
//
//func InitViper() {
//	viper.SetConfigName(os.Getenv("CONFIG"))
//	viper.SetConfigType("json")
//	viper.AddConfigPath("./config/")
//	viper.AddConfigPath(".")
//	err := viper.ReadInConfig()
//	if err != nil {
//		log.Panic("can't read config file", zap.Error(err))
//	}
//
//	err = viper.BindEnv("token")
//	if err != nil {
//		log.Panic("can't bind token env", zap.Error(err))
//	}
//	err = viper.BindEnv("admin")
//	if err != nil {
//		log.Panic("can't bind admin env", zap.Error(err))
//	}
//	err = viper.BindEnv("store_channel")
//	if err != nil {
//		log.Panic("can't bind store_channel env", zap.Error(err))
//	}
//
//	if pflag.NArg() == 1 {
//		viper.Set("token", pflag.Arg(0))
//	}
//}
//
//func main() {
//	InitViper()
//	log.InitLogger()
//	nintendo.InitClient()
//
//	botConfig := botutil.BotConfig{
//		UserProxy: viper.GetBool("bot.useProxy"),
//		ProxyUrl:  viper.GetString("bot.proxyUrl"),
//		Token:     viper.GetString("token"),
//		Debug:     viper.GetBool("bot.debug"),
//	}
//	worker := viper.GetInt("bot.worker")
//
//	myBot := botutil.NewBot(botConfig)
//	todo.InitService(myBot)
//
//	router := botutil.NewUpdateRouter()
//	router.AddCommandHandler("start", todo.Start, "start Command")
//	router.AddCommandHandler("settings", todo.Settings, "Settings Command")
//	router.AddCommandHandler("salmon_schedules", salmon.QuerySalmonSchedules, "SalmonSchedules Command")
//	router.AddCommandHandler("stages", stage.QueryStageSchedules, "Stages Command")
//	router.AddCommandHandler("help", todo.Help, "Help Command")
//	router.AddCommandHandler("help_stages", todo.HelpStages, "Help Stages Command")
//	router.AddCallbackQueryHandler(todo.ReturnToSettingsKeyboardPrefix, todo.ReturnToSettings, "Return to Settings Callback")
//	router.AddCallbackQueryHandler(todo.AccountSettingsKeyboardPrefix, todo.AccountSetting, "Settings Callback")
//	router.AddCallbackQueryHandler(todo.AccountSettingsAddKeyboardPrefix, todo.InsertAccount, "Settings Callback")
//	router.AddCallbackQueryHandler(todo.AccountSettingsDeleteKeyboardPrefix, todo.DeleteAccount, "Settings Callback")
//	router.AddCallbackQueryHandler(todo.LanguageSettingsKeyboardPrefix, todo.SetLanguage, "Settings Callback")
//	router.AddCallbackQueryHandler(todo.TimezoneSettingsKeyboardPrefix, todo.SetTimezone, "Settings Callback")
//	router.AddCallbackQueryHandler(todo.LanguageSelectionKeyboardPrefix, todo.SelectLanguage, "Select Language Callback")
//	router.AddCallbackQueryHandler(todo.TimezoneSelectionKeyboardPrefix, todo.SelectTimezone, "Select Language Callback")
//	router.SetTextHandler(todo.InputRedirectLink, "Input Redirect Link")
//
//	updateConfig := botApi.UpdateConfig{Offset: 0, Timeout: 60}
//
//	botutil.RunBotInPullMode(myBot, router, updateConfig, worker)
//}

import (
	"telegram-splatoon2-bot/app"
)

func main() {
	app.TelegramApp()
}