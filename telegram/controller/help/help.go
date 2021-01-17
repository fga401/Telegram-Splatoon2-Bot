package help

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	botMessage "telegram-splatoon2-bot/telegram/controller/internal/message"
)

func (ctrl *helpCtrl) help(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	msg := getHelpMessage(ctrl.languageSvc.Printer(status.Language), update)
	_, err := ctrl.bot.Send(msg)
	return err
}

func getHelpMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyHelp)
	msg := botMessage.NewByUpdate(update, text, nil)
	return msg
}

func (ctrl *helpCtrl) helpStages(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	msg := getHelpStageSchedulesMessage(ctrl.languageSvc.Printer(status.Language), update)
	_, err := ctrl.bot.Send(msg)
	return err
}

func getHelpStageSchedulesMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyHelpStageSchedules)
	msg := botMessage.NewByUpdate(update, text, nil)
	return msg
}
