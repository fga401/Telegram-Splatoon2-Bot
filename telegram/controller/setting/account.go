package setting

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"golang.org/x/text/message"
	userSvc "telegram-splatoon2-bot/service/user"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	messageUtil "telegram-splatoon2-bot/telegram/controller/internal/messageutil"
)

func (ctrl *settingsCtrl) accountSetting(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	accounts, err := ctrl.userSvc.ListAccounts(status.UserID)
	if err != nil {
		return errors.Wrap(err, "can't fetch accounts from database")
	}
	if len(accounts) == 0 {
		return ctrl.accountManager(update, argManager, args...)
	}
	msg := getAccountSettingMessage(ctrl.languageSvc.Printer(status.Language), update, accounts, currentAccount(accounts, status))
	_, err = ctrl.bot.Send(msg)
	return err
}

func currentAccount(accounts []userSvc.Account, status userSvc.Status) int {
	for i, account := range accounts {
		if account.SessionToken == status.SessionToken {
			return i
		}
	}
	return -1
}

func (ctrl *settingsCtrl) accountManager(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	accounts, err := ctrl.userSvc.ListAccounts(status.UserID)
	if err != nil {
		return errors.Wrap(err, "can't fetch accounts")
	}
	permission, err := ctrl.userSvc.GetPermission(status.UserID)
	if err != nil {
		return errors.Wrap(err, "can't fetch permission")
	}
	addable := permission.MaxAccount > int32(len(accounts))
	deletable := len(accounts) > 0
	msg := getAccountManagerMessage(ctrl.languageSvc.Printer(status.Language), update, accounts, permission, currentAccount(accounts, status), addable, deletable)
	_, err = ctrl.bot.Send(msg)
	return err
}

func (ctrl *settingsCtrl) accountSwitch(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	tagArgIdx := argManager.Index(ctrl.callbackQueryAdapter)[0]
	tag := args[tagArgIdx].(string)
	err := ctrl.userSvc.SwitchAccount(status.UserID, tag)
	if err != nil {
		return errors.Wrap(err, "can't switch account")
	}
	msg := getAccountSwitchMessage(ctrl.languageSvc.Printer(status.Language), update, tag)
	_, err = ctrl.bot.Send(msg)
	return err
}

func (ctrl *settingsCtrl) accountAddition(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	accounts, err := ctrl.userSvc.ListAccounts(status.UserID)
	if err != nil {
		return errors.Wrap(err, "can't fetch accounts from database")
	}
	permission, err := ctrl.userSvc.GetPermission(status.UserID)
	if err != nil {
		return errors.Wrap(err, "can't fetch permission from database")
	}
	addable := permission.MaxAccount > int32(len(accounts))
	deletable := len(accounts) > 0
	msg := getAccountManagerMessage(ctrl.languageSvc.Printer(status.Language), update, accounts, permission, currentAccount(accounts, status), addable, deletable)
	_, err = ctrl.bot.Send(msg)
	return err
}

func (ctrl *settingsCtrl) accountDeletion(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	accounts, err := ctrl.userSvc.ListAccounts(status.UserID)
	if err != nil {
		return errors.Wrap(err, "can't fetch accounts from database")
	}
	permission, err := ctrl.userSvc.GetPermission(status.UserID)
	if err != nil {
		return errors.Wrap(err, "can't fetch permission from database")
	}
	addable := permission.MaxAccount > int32(len(accounts))
	deletable := len(accounts) > 0
	msg := getAccountManagerMessage(ctrl.languageSvc.Printer(status.Language), update, accounts, permission, currentAccount(accounts, status), addable, deletable)
	_, err = ctrl.bot.Send(msg)
	return err
}

const (
	textKeyAccountSetting            = "Switch or manage your accounts:"
	textKeyAccountTagKeyboard        = "   %s"
	textKeyAccountCurrentTagKeyboard = " * %s"
	textKeyAccountManagerKeyboard    = "Manage Accounts"
)

var accountSettingMarkup = func(printer *message.Printer, accounts []userSvc.Account, current int) botApi.InlineKeyboardMarkup {
	list := make([][]botApi.InlineKeyboardButton, 0)
	for i, account := range accounts {
		textKey := textKeyAccountTagKeyboard
		if i == current {
			textKey = textKeyAccountCurrentTagKeyboard
		}
		list = append(list,
			botApi.NewInlineKeyboardRow(
				botApi.NewInlineKeyboardButtonData(
					printer.Sprintf(textKey, account.Tag),
					callbackQueryUtil.SetPrefix(KeyboardPrefixAccountSwitch, account.Tag),
				),
			),
		)
	}
	// append manager
	list = append(list,
		botApi.NewInlineKeyboardRow(
			botApi.NewInlineKeyboardButtonData(
				printer.Sprintf(textKeyAccountManagerKeyboard),
				callbackQueryUtil.SetPrefix(KeyboardPrefixAccountManager, ""),
			),
		),
	)
	ret := botApi.InlineKeyboardMarkup{
		InlineKeyboard: list,
	}
	return messageUtil.AppendBackButton(ret, KeyboardPrefixReturnToSetting, printer)
}

func getAccountSettingMessage(printer *message.Printer, update botApi.Update, accounts []userSvc.Account, current int) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountSetting)
	editMsg := update.CallbackQuery.Message
	msg := botApi.NewEditMessageText(editMsg.Chat.ID, editMsg.MessageID, text)
	markup := accountSettingMarkup(printer, accounts, current)
	msg.ReplyMarkup = &markup
	msg.ParseMode = "Markdown"
	return msg
}

const (
	textKeyAccountAddition         = "You have no account now *(%d/%d)*. You can add a new account:"
	textKeyAccountDeletion         = "Your number of account reaches the limitation *(%d/%d)*. You can delete some accounts:"
	textKeyAccountAddOrDel         = "Here are your accounts *(%d/%d)*. You can delete some accounts or add a new account:"
	textKeyAccountAdditionKeyboard = "Add Account"
)

var accountManagerMarkup = func(printer *message.Printer, accounts []userSvc.Account, current int) botApi.InlineKeyboardMarkup {
	list := make([][]botApi.InlineKeyboardButton, 0)
	for i, account := range accounts {
		textKey := textKeyAccountTagKeyboard
		if i == current {
			textKey = textKeyAccountCurrentTagKeyboard
		}
		list = append(list,
			botApi.NewInlineKeyboardRow(
				botApi.NewInlineKeyboardButtonData(
					printer.Sprintf(textKey, account.Tag),
					callbackQueryUtil.SetPrefix(KeyboardPrefixAccountDeletion, account.Tag),
				),
			),
		)
	}
	// append add
	list = append(list,
		botApi.NewInlineKeyboardRow(
			botApi.NewInlineKeyboardButtonData(
				printer.Sprintf(textKeyAccountAdditionKeyboard),
				callbackQueryUtil.SetPrefix(KeyboardPrefixAccountAddition, ""),
			),
		),
	)
	ret := botApi.InlineKeyboardMarkup{
		InlineKeyboard: list,
	}
	return messageUtil.AppendBackButton(ret, KeyboardPrefixReturnToSetting, printer)
}

func getAccountManagerMessage(printer *message.Printer, update botApi.Update, accounts []userSvc.Account, permission userSvc.Permission, current int, addable, deletable bool) botApi.Chattable {
	var text string
	if addable && deletable {
		text = printer.Sprintf(textKeyAccountAddOrDel, len(accounts), permission.MaxAccount)
	} else if addable && !deletable {
		text = printer.Sprintf(textKeyAccountAddition, len(accounts), permission.MaxAccount)
	} else if !addable && deletable {
		text = printer.Sprintf(textKeyAccountDeletion, len(accounts), permission.MaxAccount)
	}
	editMsg := update.CallbackQuery.Message
	msg := botApi.NewEditMessageText(editMsg.Chat.ID, editMsg.MessageID, text)
	markup := accountManagerMarkup(printer, accounts, current)
	msg.ReplyMarkup = &markup
	msg.ParseMode = "Markdown"
	return msg
}

const (
	textKeyAccountSwitchSuccess = "Switch your account to *%s* successfully. Use /settings to change other settings."
)

func getAccountSwitchMessage(printer *message.Printer, update botApi.Update, tag string) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountSwitchSuccess, tag)
	editMsg := update.CallbackQuery.Message
	msg := botApi.NewEditMessageText(editMsg.Chat.ID, editMsg.MessageID, text)
	msg.ParseMode = "Markdown"
	return msg
}
