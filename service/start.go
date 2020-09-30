package service

import (
	"fmt"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"strconv"
	"telegram-splatoon2-bot/common/cache"
	"telegram-splatoon2-bot/nintendo"
)

func Start(update *botapi.Update, bot *botapi.BotAPI) error {
	msg := update.Message
	replyMsg := botapi.NewMessage(update.Message.Chat.ID, "update.Message.Text")
	replyMsg.ReplyToMessageID = update.Message.MessageID
	proofKey, err := nintendo.NewProofKey()
	if err != nil {
		return errors.Wrap(err, "can't generate proof key")
	}
	link, err := nintendo.NewLoginLink(proofKey)
	if err != nil {
		return errors.Wrap(err, "can't generate login link")
	}
	cache.GoCache.SetDefault(getProofKeyCacheKey(msg.From), proofKey)
	replyMsg.Text = link
	_, err = bot.Send(replyMsg)
	if err != nil {
		return errors.Wrap(err, "can't send message")
	}
	return nil
}

func InputRedirectLink(update *botapi.Update, bot *botapi.BotAPI) error {
	text := update.Message.Text
	user := update.Message.From
	proofKeyInterface, in := cache.GoCache.Get(getProofKeyCacheKey(user))
	if !in {
		return fmt.Errorf("unknown input")
	}
	proofKey := proofKeyInterface.([]byte)
	cookies, err := nintendo.GetCookies(text, proofKey)
	if err != nil {
		return errors.Wrap(err, "login failed")
	}
	// todo: temporary test
	replyMsg := botapi.NewMessage(update.Message.Chat.ID, "update.Message.Text")
	replyMsg.ReplyToMessageID = update.Message.MessageID
	replyMsg.Text = cookies
	_, err = bot.Send(replyMsg)
	if err != nil {
		return errors.Wrap(err, "can't send message")
	}
	return nil
}

func register(user *botapi.User) error {
	if user.UserName == "" {
		return fmt.Errorf("UserName not existed")
	}

	return nil
}

func getProofKeyCacheKey(user *botapi.User) string {
	return "ProofKey_" + strconv.FormatInt(int64(user.ID), 10)
}

