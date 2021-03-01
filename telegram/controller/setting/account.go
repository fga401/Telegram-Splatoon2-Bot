package setting

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/nintendo"
	userSvc "telegram-splatoon2-bot/service/user"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	"telegram-splatoon2-bot/telegram/controller/internal/convert"
	"telegram-splatoon2-bot/telegram/controller/internal/markup"
	botMessage "telegram-splatoon2-bot/telegram/controller/internal/message"
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
	msg := getAccountSwitchMessage(ctrl.languageSvc.Printer(status.Language), update)
	_, _ = ctrl.bot.Send(msg)
	err := ctrl.userSvc.SwitchAccount(status.UserID, tag)
	if err != nil {
		return errors.Wrap(err, "can't switch account")
	}
	msg = getAccountSwitchSuccessMessage(ctrl.languageSvc.Printer(status.Language), update, tag)
	_, _ = ctrl.bot.Send(msg)
	return convert.CallbackQueryToCommand(ctrl.AccountSetting)(update)
}

func (ctrl *settingsCtrl) accountAddition(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	loginLink, err := ctrl.userSvc.NewLoginLink(status.UserID)
	if err != nil {
		return errors.Wrap(err, "can't generate new login link")
	}
	msg := getAccountAdditionMessage(ctrl.languageSvc.Printer(status.Language), update, loginLink)
	_, err = ctrl.bot.Send(msg)
	return err
}

func (ctrl *settingsCtrl) accountRedirectLink(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	redirectLink := update.Message.Text
	if !nintendo.IsRedirectLinkValid(redirectLink) {
		msg := getAccountRedirectLinkInvalidMessage(ctrl.languageSvc.Printer(status.Language), update)
		_, err := ctrl.bot.Send(msg)
		return err
	}

	msg := getAccountRedirectLinkFetchingMessage(ctrl.languageSvc.Printer(status.Language), update)
	resp, err := ctrl.bot.Send(msg)
	if err != nil {
		return err
	}
	account, err := ctrl.userSvc.AddAccount(status.UserID, redirectLink)
	if errors.Is(err, userSvc.ErrNoProofKey{}) {
		msg := getAccountRedirectLinkNoProofKeyMessage(ctrl.languageSvc.Printer(status.Language), resp)
		_, _ = ctrl.bot.Send(msg)
		return ctrl.AccountManager(update)
	}
	if errors.Is(err, userSvc.ErrAccountExisted{}) {
		msg := getAccountRedirectLinkAccountExistedMessage(ctrl.languageSvc.Printer(status.Language), resp)
		_, _ = ctrl.bot.Send(msg)
		return ctrl.AccountManager(update)
	}
	if err != nil {
		log.Error("internal error", zap.Error(err))
		msg := getAccountRedirectLinkOtherErrorMessage(ctrl.languageSvc.Printer(status.Language), resp)
		_, err = ctrl.bot.Send(msg)
		return err
	}
	msg = getAccountRedirectLinkSuccessMessage(ctrl.languageSvc.Printer(status.Language), resp, account.Tag)
	_, _ = ctrl.bot.Send(msg)
	return ctrl.AccountManager(update)
}

func (ctrl *settingsCtrl) accountDeletionConfirm(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	tagArgIdx := argManager.Index(ctrl.callbackQueryAdapter)[0]
	tag := args[tagArgIdx].(string)
	msg := getAccountDeletionConfirmMessage(ctrl.languageSvc.Printer(status.Language), update, tag)
	_, err := ctrl.bot.Send(msg)
	return err
}

func (ctrl *settingsCtrl) accountDeletion(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	tagArgIdx := argManager.Index(ctrl.callbackQueryAdapter)[0]
	tag := args[tagArgIdx].(string)
	msg := getAccountDeletionMessage(ctrl.languageSvc.Printer(status.Language), update)
	_, _ = ctrl.bot.Send(msg)
	err := ctrl.userSvc.DeleteAccount(status.UserID, tag)
	if err != nil {
		return errors.Wrap(err, "can't delete account")
	}
	msg = getAccountDeletionSuccessMessage(ctrl.languageSvc.Printer(status.Language), update, tag)
	_, _ = ctrl.bot.Send(msg)
	return convert.CallbackQueryToCommand(ctrl.AccountManager)(update)
}

const (
	textKeyAccountSetting                    = "Switch or manage your accounts:"
	textKeyAccountTagKeyboard                = "   %s"
	textKeyAccountCurrentTagKeyboard         = " * %s"
	textKeyAccountDeletionTagKeyboard        = "   Delete %s"
	textKeyAccountDeletionCurrentTagKeyboard = " * Delete %s"
	textKeyAccountManagerKeyboard            = "Manage Accounts"
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
				printer.Sprintf(textKeyAccountManagerKeyboard, len(accounts)),
				callbackQueryUtil.SetPrefix(KeyboardPrefixAccountManager, ""),
			),
		),
	)
	ret := botApi.InlineKeyboardMarkup{
		InlineKeyboard: list,
	}
	return markup.AppendBackButton(ret, KeyboardPrefixSetting, printer)
}

func getAccountSettingMessage(printer *message.Printer, update botApi.Update, accounts []userSvc.Account, current int) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountSetting, len(accounts))
	markup := accountSettingMarkup(printer, accounts, current)
	msg := botMessage.NewByUpdate(update, text, &markup)
	return msg
}

const (
	textKeyAccountManagerAddition         = "You have no account now *(%d/%d)*. You can add a new account:"
	textKeyAccountManagerDeletion         = "Your number of account reaches the limitation *(%d/%d)*. You can delete some accounts:"
	textKeyAccountManagerAddOrDel         = "Here are your accounts *(%d/%d)*. You can delete some accounts or add a new account:"
	textKeyAccountManagerAdditionKeyboard = "Add Account"
)

var accountManagerMarkup = func(printer *message.Printer, accounts []userSvc.Account, current int, addable, deletable bool) botApi.InlineKeyboardMarkup {
	list := make([][]botApi.InlineKeyboardButton, 0)
	for i, account := range accounts {
		textKey := textKeyAccountDeletionTagKeyboard
		if i == current {
			textKey = textKeyAccountDeletionCurrentTagKeyboard
		}
		list = append(list,
			botApi.NewInlineKeyboardRow(
				botApi.NewInlineKeyboardButtonData(
					printer.Sprintf(textKey, account.Tag),
					callbackQueryUtil.SetPrefix(KeyboardPrefixAccountDeletionConfirm, account.Tag),
				),
			),
		)
	}
	// append add
	if addable {
		list = append(list,
			botApi.NewInlineKeyboardRow(
				botApi.NewInlineKeyboardButtonData(
					printer.Sprintf(textKeyAccountManagerAdditionKeyboard),
					callbackQueryUtil.SetPrefix(KeyboardPrefixAccountAddition, ""),
				),
			),
		)
	}
	ret := botApi.InlineKeyboardMarkup{
		InlineKeyboard: list,
	}
	if deletable {
		ret = markup.AppendBackButton(ret, KeyboardPrefixAccountSetting, printer)
	} else {
		ret = markup.AppendBackButton(ret, KeyboardPrefixSetting, printer)
	}
	return ret
}

func getAccountManagerMessage(printer *message.Printer, update botApi.Update, accounts []userSvc.Account, permission userSvc.Permission, current int, addable, deletable bool) botApi.Chattable {
	var text string
	if addable && deletable {
		text = printer.Sprintf(textKeyAccountManagerAddOrDel, len(accounts), permission.MaxAccount)
	} else if addable && !deletable {
		text = printer.Sprintf(textKeyAccountManagerAddition, len(accounts), permission.MaxAccount)
	} else if !addable && deletable {
		text = printer.Sprintf(textKeyAccountManagerDeletion, len(accounts), permission.MaxAccount)
	}
	markup := accountManagerMarkup(printer, accounts, current, addable, deletable)
	msg := botMessage.NewByUpdate(update, text, &markup)
	return msg
}

const (
	textKeyAccountSwitch        = "Switching account..."
	textKeyAccountSwitchSuccess = "Switch your account to *%s* successfully."
)

func getAccountSwitchMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountSwitch)
	msg := botMessage.NewByUpdate(update, text, nil)
	return msg
}

func getAccountSwitchSuccessMessage(printer *message.Printer, update botApi.Update, tag string) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountSwitchSuccess, tag)
	msg := botMessage.NewByUpdate(update, text, nil)
	return msg
}

const (
	textKeyAccountAddition         = "Please open the following link and:\n*1.* Login;\n*2.* Right click / Long press the *<Select this account>* button;\n*3.* Copy the link address;\n*4.* Paste and send to this bot."
	textKeyAccountAdditionKeyboard = "Login Link"
)

func getAccountAdditionMessage(printer *message.Printer, update botApi.Update, link string) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountAddition)
	markup := botApi.NewInlineKeyboardMarkup(botApi.NewInlineKeyboardRow(botApi.NewInlineKeyboardButtonURL(textKeyAccountAdditionKeyboard, link)))
	msg := botMessage.NewByUpdate(update, text, &markup)
	return msg
}

const (
	textKeyAccountRedirectLinkInvalid        = "Invalid input. Please use /help to get usage."
	textKeyAccountRedirectLinkFetching       = "Fetching account from Nintendo server..."
	textKeyAccountRedirectLinkNoProofKey     = "Your link is expired. Please use *Add* *Account* to regenerate."
	textKeyAccountRedirectLinkAccountExisted = "Your account is already existed. Please use *Add* *Account* to add another one."
	textKeyAccountRedirectLinkOtherError     = "Internal error. Please paste your link and retry."
	textKeyAccountRedirectLinkSuccess        = "Account *%s* has been added."
)

func getAccountRedirectLinkInvalidMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountRedirectLinkInvalid)
	msg := botMessage.NewByUpdate(update, text, nil)
	return msg
}

func getAccountRedirectLinkFetchingMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountRedirectLinkFetching)
	msg := botMessage.NewByUpdate(update, text, nil)
	return msg
}

func getAccountRedirectLinkNoProofKeyMessage(printer *message.Printer, msg *botApi.Message) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountRedirectLinkNoProofKey)
	ret := botMessage.NewByMsg(msg, text, nil, true)
	return ret
}

func getAccountRedirectLinkAccountExistedMessage(printer *message.Printer, msg *botApi.Message) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountRedirectLinkAccountExisted)
	ret := botMessage.NewByMsg(msg, text, nil, true)
	return ret
}

func getAccountRedirectLinkOtherErrorMessage(printer *message.Printer, msg *botApi.Message) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountRedirectLinkOtherError)
	ret := botMessage.NewByMsg(msg, text, nil, true)
	return ret
}

func getAccountRedirectLinkSuccessMessage(printer *message.Printer, msg *botApi.Message, tag string) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountRedirectLinkSuccess, tag)
	ret := botMessage.NewByMsg(msg, text, nil, true)
	return ret
}

const (
	textKeyAccountDeletionConfirm = "Are you sure to delete this account?\n*%s*"

	textKeyYes = "Yes"
	textKeyNo  = "No"
)

var accountDeletionConfirmMarkup = func(printer *message.Printer, tag string) botApi.InlineKeyboardMarkup {
	list := make([][]botApi.InlineKeyboardButton, 0)
	list = append(list,
		botApi.NewInlineKeyboardRow(
			botApi.NewInlineKeyboardButtonData(
				printer.Sprintf(textKeyNo),
				callbackQueryUtil.SetPrefix(KeyboardPrefixAccountManager, ""),
			),
			botApi.NewInlineKeyboardButtonData(
				printer.Sprintf(textKeyYes),
				callbackQueryUtil.SetPrefix(KeyboardPrefixAccountDeletion, tag),
			),
		),
	)
	ret := botApi.InlineKeyboardMarkup{
		InlineKeyboard: list,
	}
	return ret
}

func getAccountDeletionConfirmMessage(printer *message.Printer, update botApi.Update, tag string) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountDeletionConfirm, tag)
	markup := accountDeletionConfirmMarkup(printer, tag)
	msg := botMessage.NewByUpdate(update, text, &markup)
	return msg
}

const (
	textKeyAccountDeletion        = "Deleting account..."
	textKeyAccountDeletionSuccess = "Delete account *%s* successfully."
)

func getAccountDeletionMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountDeletion)
	msg := botMessage.NewByUpdate(update, text, nil)
	return msg
}

func getAccountDeletionSuccessMessage(printer *message.Printer, update botApi.Update, tag string) botApi.Chattable {
	text := printer.Sprintf(textKeyAccountDeletionSuccess, tag)
	msg := botMessage.NewByUpdate(update, text, nil)
	return msg
}
