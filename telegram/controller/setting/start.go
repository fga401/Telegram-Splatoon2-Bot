package setting

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/message"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/language"
	userSvc "telegram-splatoon2-bot/service/user"
)

func (ctrl *settingsCtrl) Start(update botApi.Update) error {
	user := update.Message.From
	userID := userSvc.ID(user.ID)
	existed, err := ctrl.userSvc.Existed(userID)
	if err != nil {
		return errors.Wrap(err, "can't check if user existed")
	}
	if !existed {
		err = ctrl.userSvc.Register(userID, user.UserName)
		if err != nil {
			return errors.Wrap(err, "can't register new user")
		}
		log.Info("new user register", zap.Object("user", log.UserPtrLogger(user)))
		msg := getStartMessage(ctrl.languageSvc.Printer(language.English), update)
		_, err = ctrl.bot.Send(msg)
		if err != nil {
			log.Warn("can't send hello message", zap.Object("update", log.UpdateLogger(update)), zap.Error(err))
		}
	}
	return ctrl.Setting(update)
}

const (
	textKeyWelcome = "Welcome to use this bot."
)

func getStartMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyWelcome)
	msg := botApi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	return msg
}
