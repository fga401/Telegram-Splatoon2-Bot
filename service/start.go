package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common"
	"telegram-splatoon2-bot/nintendo"
)

func Start(update *botapi.Update, bot *botapi.BotAPI) {
	err := start(update, bot)
	if err != nil {
		logger.Error("Start Command",
			zap.Bool("status", false),
			zap.Object("message", common.WrapMessage(update.Message)),
			zap.Error(err))
	} else {
		logger.Info("Start Command",
			zap.Bool("status", true),
			zap.Object("message", common.WrapMessage(update.Message)))
	}
}

func start(update *botapi.Update, bot *botapi.BotAPI) error {
	msg := botapi.NewMessage(update.Message.Chat.ID, "update.Message.Text")
	msg.ReplyToMessageID = update.Message.MessageID
	link, err := nintendo.GetLoginLink()
	if err != nil {
		return errors.Wrap(err, "can't get login link")
	}
	msg.Text = link
	_, err = bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "can't send message")
	}
	return nil
}