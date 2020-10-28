package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func Help(update *botapi.Update) error {
	user := update.Message.From
	runtime, err := FetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	texts := getI18nText(runtime.Language, user, NewI18nKey(helpTextKey))
	msg := botapi.NewMessage(update.Message.Chat.ID, texts[0])
	msg.ParseMode = "Markdown"
	err = sendWithRetry(bot, msg)
	return err
}

func HelpStages(update *botapi.Update) error {
	user := update.Message.From
	runtime, err := FetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	texts := getI18nText(runtime.Language, user, NewI18nKey(helpStagesTextKey))
	msg := botapi.NewMessage(update.Message.Chat.ID, texts[0])
	msg.ParseMode = "Markdown"
	err = sendWithRetry(bot, msg)
	return err
}
