package setting

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/language"
	userSvc "telegram-splatoon2-bot/service/user"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	"telegram-splatoon2-bot/telegram/controller/internal/convert"
	"telegram-splatoon2-bot/telegram/controller/internal/markup"
	botMessage "telegram-splatoon2-bot/telegram/controller/internal/message"
)

func (ctrl *settingsCtrl) languageSetting(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	msg := getLanguageSettingMessage(ctrl.languageSvc.Printer(status.Language), update, ctrl.languageSvc)
	_, err := ctrl.bot.Send(msg)
	return err
}

func (ctrl *settingsCtrl) languageSelection(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	ietfIdx := argManager.Index(ctrl.callbackQueryAdapter)[0]
	ietf := args[ietfIdx].(string)
	lang := language.ByIETF(ietf)
	status := args[statusArgIdx].(userSvc.Status)
	err := ctrl.userSvc.UpdateStatusLanguage(status.UserID, lang)
	if err != nil {
		return errors.Wrap(err, "can't update language")
	}
	status.Language = lang
	log.Info("user language updated",
		zap.String("language", ietf),
		zap.Object("user", log.UserPtrLogger(update.CallbackQuery.From)),
	)
	msg := getLanguageSelectionMessage(ctrl.languageSvc.Printer(status.Language), update, lang)
	_, err = ctrl.bot.Send(msg)
	return convert.CallbackQueryToCommand(ctrl.Setting)(update)
}

const (
	textKeyLanguageSelection        = "Please select your preferred language:"
	textKeyLanguageSelectionSuccess = "Change your language to *%s* successfully."
)

func languageKey(lang language.Language) string {
	return "lang: " + lang.IETF()
}

var languageSettingMarkup = func(printer *message.Printer, languageSvc language.Service) botApi.InlineKeyboardMarkup {
	list := make([][]botApi.InlineKeyboardButton, 0)
	for _, lang := range languageSvc.Supported() {
		list = append(list,
			botApi.NewInlineKeyboardRow(
				botApi.NewInlineKeyboardButtonData(
					printer.Sprintf(languageKey(lang)),
					callbackQueryUtil.SetPrefix(KeyboardPrefixLanguageSelection, lang.IETF()),
				),
			),
		)
	}
	ret := botApi.InlineKeyboardMarkup{
		InlineKeyboard: list,
	}
	return markup.AppendBackButton(ret, KeyboardPrefixSetting, printer)
}

func getLanguageSettingMessage(printer *message.Printer, update botApi.Update, langSvc language.Service) botApi.Chattable {
	text := printer.Sprintf(textKeyLanguageSelection)
	markup := languageSettingMarkup(printer, langSvc)
	msg := botMessage.NewByUpdate(update, text, &markup)
	return msg
}

func getLanguageSelectionMessage(printer *message.Printer, update botApi.Update, lang language.Language) botApi.Chattable {
	text := printer.Sprintf(textKeyLanguageSelectionSuccess, printer.Sprintf(languageKey(lang)))
	msg := botMessage.NewByUpdate(update, text, nil)
	return msg
}
