package botutil

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"net/http"
	proxy2 "telegram-splatoon2-bot/common/proxy"
	log "telegram-splatoon2-bot/logger"
)

type Job struct {
	bot              *botapi.BotAPI
	update           *botapi.Update
	describedHandler *DescribedHandler
}

type BotConfig struct {
	UserProxy bool
	ProxyUrl  string
	Token     string
	Debug     bool
}

func NewBot(config BotConfig) *botapi.BotAPI {
	useProxy := config.UserProxy
	proxyUrl := config.ProxyUrl
	debug := config.Debug
	token := config.Token
	proxy := proxy2.GetProxy()
	if proxyUrl != "" {
		proxy = proxy2.GetProxyWithUrl(proxyUrl)
	}
	if !useProxy {
		proxy = nil
	}
	client := &http.Client{Transport: &http.Transport{Proxy: proxy}}

	bot, err := botapi.NewBotAPIWithClient(token, client)
	if err != nil {
		log.Fatal("Bot API initialization failed.", zap.Error(err))
	}
	bot.Debug = debug
	log.Info("Authorized on account.", zap.String("account", bot.Self.UserName))
	return bot
}

func RunBotInPullMode(bot *botapi.BotAPI, router *UpdateRouter, updateConfig botapi.UpdateConfig, worker int) {
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal("can't get bot update channel", zap.Error(err))
	}
	jobChan := make(chan Job, worker)
	for i := 0; i < worker; i++ {
		go func() {
			for job := range jobChan {
				err := job.describedHandler.handler(job.update, job.bot)
				if err != nil {
					log.Error(job.describedHandler.des,
						zap.Bool("status", false),
						zap.Object("update", log.WrapUpdate(job.update)),
						zap.Error(err))
				} else {
					log.Info(job.describedHandler.des,
						zap.Bool("status", true),
						zap.Object("update", log.WrapUpdate(job.update)))
				}
			}
		}()
	}

	for update := range updates {
		describedHandler := router.Route(&update)
		if describedHandler != nil {
			jobChan <- Job{
				bot:              bot,
				update:           &update,
				describedHandler: describedHandler,
			}
		}
	}
}
