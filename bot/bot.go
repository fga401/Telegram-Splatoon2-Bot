package bot

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"telegram-splatoon2-bot/common"
	log "telegram-splatoon2-bot/logger"
)

func NewBot(token string) *botapi.BotAPI{
	useProxy := viper.GetBool("bot.useProxy")
	proxy := common.GetProxy()
	if viper.InConfig("bot.proxyUrl"){
		proxy = common.GetProxyWithUrl(viper.GetString("nintendo.proxyUrl"))
	}
	if !useProxy {
		proxy = nil
	}
	client := &http.Client{Transport: &http.Transport{Proxy: proxy}}

	bot, err := botapi.NewBotAPIWithClient(token, client)
	if err != nil {
		log.Fatal("Bot API initialization failed.", zap.Error(err))
	}
	bot.Debug = true
	log.Info("Authorized on account.", zap.String("account", bot.Self.UserName))
	return bot
}

func RunBotInPullMode(bot *botapi.BotAPI, router *CommandRouter, updateConfig botapi.UpdateConfig){
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal("can't get bot update channel", zap.Error(err))
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Info("message received",
			zap.Object("message", log.WrapMessage(update.Message)))

		if update.Message.IsCommand() {
			router.Run(&update, bot)
		}
	}
}