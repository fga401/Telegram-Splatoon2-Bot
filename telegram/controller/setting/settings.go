package setting

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
	userSvc "telegram-splatoon2-bot/service/user"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
)

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
	return botApi.NewInlineKeyboardMarkup(
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
}

func getMainSettingMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeySetting)
	markup := mainSettingMarkup(printer)
	msg := botApi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = markup
	msg.ParseMode = "Markdown"
	return msg
}
