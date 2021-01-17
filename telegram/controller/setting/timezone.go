package setting

import (
	"strconv"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/timezone"
	userSvc "telegram-splatoon2-bot/service/user"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	messageUtil "telegram-splatoon2-bot/telegram/controller/internal/messageutil"
)

func (ctrl *settingsCtrl) timezoneSetting(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	msg := getTimezoneSettingMessage(ctrl.languageSvc.Printer(status.Language), update)
	_, err := ctrl.bot.Send(msg)
	return err
}

func (ctrl *settingsCtrl) timezoneSelection(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	tzIdx := argManager.Index(ctrl.callbackQueryAdapter)[0]
	tzInMin, _ := strconv.Atoi(args[tzIdx].(string))
	tz := timezone.ByMinute(tzInMin)
	status := args[statusArgIdx].(userSvc.Status)
	err := ctrl.userSvc.UpdateStatusTimezone(status.UserID, tz)
	if err != nil {
		return errors.Wrap(err, "can't update timezone")
	}
	status.Timezone = tz
	log.Info("user timezone updated",
		zap.Int("timezone", tzInMin),
		zap.Object("user", log.UserPtrLogger(update.CallbackQuery.From)),
	)
	msg := getTimezoneSelectionMessage(ctrl.languageSvc.Printer(status.Language), update, tz)
	_, err = ctrl.bot.Send(msg)
	return err
}

const (
	textKeyTimezoneSelection        = "Please select your timezone:"
	textKeyTimezoneSelectionSuccess = "Change your timezone to *%s* successfully! Use /settings to change other settings."
)

func timezoneKey(t timezone.Timezone) string {
	return "local: " + strconv.Itoa(int(t))
}

var timezoneSettingMarkup = func(printer *message.Printer) botApi.InlineKeyboardMarkup {
	list := make([][]botApi.InlineKeyboardButton, 0)
	for _, tz := range timezone.All {
		list = append(list,
			botApi.NewInlineKeyboardRow(
				botApi.NewInlineKeyboardButtonData(
					printer.Sprintf(timezoneKey(tz)),
					callbackQueryUtil.SetPrefix(KeyboardPrefixTimezoneSelection, strconv.Itoa(tz.Minute())),
				),
			),
		)
	}
	ret := botApi.InlineKeyboardMarkup{
		InlineKeyboard: list,
	}
	return messageUtil.AppendBackButton(ret, KeyboardPrefixReturnToSetting, printer)
}

func getTimezoneSettingMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyTimezoneSelection)
	markup := timezoneSettingMarkup(printer)
	editMsg := update.CallbackQuery.Message
	msg := botApi.NewEditMessageText(editMsg.Chat.ID, editMsg.MessageID, text)
	msg.ReplyMarkup = &markup
	msg.ParseMode = "Markdown"
	return msg
}

func getTimezoneSelectionMessage(printer *message.Printer, update botApi.Update, tz timezone.Timezone) botApi.Chattable {
	text := printer.Sprintf(textKeyTimezoneSelectionSuccess, printer.Sprintf(timezoneKey(tz)))
	editMsg := update.CallbackQuery.Message
	msg := botApi.NewEditMessageText(editMsg.Chat.ID, editMsg.MessageID, text)
	msg.ParseMode = "Markdown"
	return msg
}
