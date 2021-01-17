package setting

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
	userSvc "telegram-splatoon2-bot/service/user"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	"telegram-splatoon2-bot/telegram/controller/internal/markup"
	botMessage "telegram-splatoon2-bot/telegram/controller/internal/message"
)

func (ctrl *settingsCtrl) cancelSetting(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	msg := update.CallbackQuery.Message
	_, err := ctrl.bot.Send(botApi.NewDeleteMessage(msg.Chat.ID, msg.MessageID))
	return err
}

func (ctrl *settingsCtrl) setting(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	msg := getMainSettingMessage(ctrl.languageSvc.Printer(status.Language), update)
	_, err := ctrl.bot.Send(msg)
	return err
}

const (
	textKeySetting = "What do you want to change?"
)

var mainSettingMarkup = func(printer *message.Printer) botApi.InlineKeyboardMarkup {
	ret := botApi.NewInlineKeyboardMarkup(
		botApi.NewInlineKeyboardRow(
			botApi.NewInlineKeyboardButtonData(
				printer.Sprintf("Account"),
				callbackQueryUtil.SetPrefix(KeyboardPrefixAccountSetting, ""),
			),
			botApi.NewInlineKeyboardButtonData(
				printer.Sprintf("Language"),
				callbackQueryUtil.SetPrefix(KeyboardPrefixLanguageSettings, ""),
			),
			botApi.NewInlineKeyboardButtonData(
				printer.Sprintf("Timezone"),
				callbackQueryUtil.SetPrefix(KeyboardPrefixTimezoneSettings, ""),
			),
		),
	)
	return markup.AppendBackButton(ret, KeyboardPrefixCancelSetting, printer)
}

func getMainSettingMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeySetting)
	markup := mainSettingMarkup(printer)
	msg := botMessage.NewByUpdate(update, text, &markup)
	return msg
}
