package bot

import (
	"time"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/common/util"
	"telegram-splatoon2-bot/telegram/bot/internal/limit"
	sendMediaGroup "telegram-splatoon2-bot/telegram/bot/internal/send_media_group"
)

type Bot interface {
	Send(msg botApi.Chattable) (*botApi.Message, error)
	SendMediaGroup(config sendMediaGroup.Config) ([]*botApi.Message, error)
	AnswerCallbackQuery(chatID string, option ...CallbackQueryConfig) error
}

type impl struct {
	config Config
	bot    *botApi.BotAPI
}

func New(bot *botApi.BotAPI, config Config) Bot {
	return &impl{
		bot:    bot,
		config: config,
	}
}

func (s *impl) Send(msg botApi.Chattable) (*botApi.Message, error) {
	var respMsg botApi.Message
	err := util.Retry(func() error {
		var err error
		respMsg, err = s.bot.Send(msg)
		if is, sec := limit.IsTooManyRequestError(err); is {
			// todo: more info?
			log.Warn("send message blocked by telegram request limits", zap.Int("after", sec))
			time.Sleep(time.Duration(sec) * time.Second)
		}
		return err
	}, s.config.RetryTimes)
	if err != nil {
		err = errors.Wrap(err, "can't send message")
	}
	return &respMsg, err
}

func (s *impl) SendMediaGroup(config sendMediaGroup.Config) ([]*botApi.Message, error) {
	return sendMediaGroup.Do(s.bot, config, s.config.RetryTimes)
}

func (s *impl) AnswerCallbackQuery(callbackQueryID string, option ...CallbackQueryConfig) error {
	config := s.config.DefaultCallbackQueryConfig
	if len(option) > 0 {
		config = option[0]
	}
	err := util.Retry(func() error {
		var err error
		_, err = s.bot.AnswerCallbackQuery(botApi.CallbackConfig{
			CallbackQueryID: callbackQueryID,
			Text:            config.Text,
			ShowAlert:       config.ShowAlert,
			CacheTime:       config.CacheTime,
		})
		if is, sec := limit.IsTooManyRequestError(err); is {
			// todo: more info?
			log.Warn("AnswerCallbackQuery blocked by telegram request limits", zap.Int("after", sec))
			time.Sleep(time.Duration(sec) * time.Second)
		}
		return err
	}, s.config.RetryTimes)
	return err
}
