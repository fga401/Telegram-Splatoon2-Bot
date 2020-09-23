package main

import (
	"fmt"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
	_ "github.com/valyala/fasthttp"
	"os"
	"telegram-splatoon2-bot/bot"
	"telegram-splatoon2-bot/service"
)

func init() {
	viper.SetConfigFile(os.Getenv("CONFIG"))
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
	token := viper.GetString("token")
	myBot := bot.InitBot(token)

	router := bot.NewCommandRouter()
	router.Add("start", service.Start)

	updateConfig := botapi.UpdateConfig{Offset: 0, Timeout: 60}

	bot.RunBotInPullMode(myBot, router, updateConfig)
}


