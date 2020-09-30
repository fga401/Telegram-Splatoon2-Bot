package main

import (
	"fmt"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"telegram-splatoon2-bot/bot"
	"telegram-splatoon2-bot/common/cache"
	"telegram-splatoon2-bot/common/db"
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
		panic(fmt.Errorf("main: config file: %s \n", err))
	}

	err = viper.BindEnv("token")
	if err != nil {
		panic(errors.Wrap(err, "can't bind env"))
	}

	if pflag.NArg() == 1 {
		viper.Set("token", pflag.Arg(0))
	}
}

func main() {
	InitViper()
	logger.InitLogger()
	nintendo.InitClient()
	cache.InitCache()
	db.InitDatabaseInstance()

	botConfig := bot.BotConfig{
		UserProxy: viper.GetBool("bot.useProxy"),
		ProxyUrl:  viper.GetString("bot.proxyUrl"),
		Token:     viper.GetString("token"),
		Debug:     viper.GetBool("bot.debug"),
	}
	worker := viper.GetInt("bot.worker")

	myBot := bot.NewBot(botConfig)
	router := bot.NewUpdateRouter()
	router.AddCommandHandler("start", service.Start, "Start Command")
	router.SetTextHandler(service.InputRedirectLink, "Input Redirect Link")

	updateConfig := botapi.UpdateConfig{Offset: 0, Timeout: 60}

	bot.RunBotInPullMode(myBot, router, updateConfig, worker)
}
