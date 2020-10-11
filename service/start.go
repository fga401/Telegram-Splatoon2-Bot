package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"strconv"
	botutils "telegram-splatoon2-bot/botutil"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/nintendo"
	"telegram-splatoon2-bot/service/db"
)

func Start(update *botapi.Update) error {
	user := update.Message.From
	existed, err := UserTable.IsUserExisted(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't check if user existed")
	}
	if !existed {
		log.Info("new user", zap.Object("user", log.WrapUser(user)))
		err = register(user)
		if err != nil {
			return errors.Wrap(err, "can't register new user")
		}
		msg := botapi.NewMessage(update.Message.Chat.ID, "Welcome to use this bot.")
		msg.ParseMode = "Markdown"
		err = sendWithRetry(bot, msg)
		if err != nil {
			log.Warn("can't send hello message", zap.Object("user", log.WrapUser(user)), zap.Error(err))
		}
	}
	return Settings(update)
}

func Settings(update *botapi.Update) error {
	user := update.Message.From
	runtime, err := fetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	texts := getI18nText(runtime.Language, user, NewI18nKey(settingsTextKey))
	markup := getStaticMarkup(settingsKeyboard, runtime.Language)
	msg := botapi.NewMessage(update.Message.Chat.ID, texts[0])
	msg.ReplyMarkup = markup
	msg.ParseMode = "Markdown"
	return sendWithRetry(bot, msg)

}

func SetLanguage(update *botapi.Update) error {
	callback := update.CallbackQuery
	user := callback.From
	runtime, err := fetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}

	_, err = bot.AnswerCallbackQuery(botapi.CallbackConfig{
		CallbackQueryID: callback.ID,
		CacheTime:       callbackQueryCachedSecond,
	})
	if err != nil {
		return errors.Wrap(err, "can't answer callback query")
	}

	texts := getI18nText(runtime.Language, user, NewI18nKey(selectLanguageTextKey))
	msg := botapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, texts[0])
	markup := getStaticMarkup(languageKeyboard, runtime.Language)
	msg.ReplyMarkup = &markup
	msg.ParseMode = "Markdown"
	return sendWithRetry(bot, msg)
}

func SelectLanguage(update *botapi.Update) error {
	callback := update.CallbackQuery
	user := callback.From
	runtime, err := fetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	lang := botutils.GetCallbackQueryOriginText(callback.Data)
	if lang == "BACK" {
		texts := getI18nText(runtime.Language, user, NewI18nKey(settingsTextKey))
		markup := getStaticMarkup(settingsKeyboard, runtime.Language)
		msg := botapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, texts[0])
		msg.ReplyMarkup = &markup
		msg.ParseMode = "Markdown"
		return sendWithRetry(bot, msg)
	}

	if _, err = language.Parse(runtime.Language); err != nil {
		return errors.Wrap(err, "unknown language")
	}
	runtime.Language = lang
	err = RuntimeTable.UpdateRuntimeLanguage(runtime.Uid, runtime.Language)
	if err != nil {
		return errors.Wrap(err, "can't update language")
	}
	log.Info("user language updated", zap.String("language", lang), zap.Object("user", log.WrapUser(user)))
	Cache.DeleteRuntime(int64(user.ID))

	texts := getI18nText(runtime.Language, user, NewI18nKey(selectLanguageSuccessfullyTextKey, langToText(runtime.Language, runtime.Language)))
	msg := botapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, texts[0])
	msg.ParseMode = "Markdown"
	return sendWithRetry(bot, msg)
}

func SetTimezone(update *botapi.Update) error {
	callback := update.CallbackQuery
	user := callback.From
	runtime, err := fetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}

	_, err = bot.AnswerCallbackQuery(botapi.CallbackConfig{
		CallbackQueryID: callback.ID,
		CacheTime:       callbackQueryCachedSecond,
	})
	if err != nil {
		return errors.Wrap(err, "can't answer callback query")
	}

	texts := getI18nText(runtime.Language, user, NewI18nKey(selectTimezoneTextKey))
	msg := botapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, texts[0])
	markup := getStaticMarkup(timezoneKeyboard, runtime.Language)
	msg.ReplyMarkup = &markup
	msg.ParseMode = "Markdown"
	return sendWithRetry(bot, msg)
}

func SelectTimezone(update *botapi.Update) error {
	callback := update.CallbackQuery
	user := callback.From
	runtime, err := fetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	data := botutils.GetCallbackQueryOriginText(callback.Data)
	if data == "BACK" {
		texts := getI18nText(runtime.Language, user, NewI18nKey(settingsTextKey))
		markup := getStaticMarkup(settingsKeyboard, runtime.Language)
		msg := botapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, texts[0])
		msg.ReplyMarkup = &markup
		msg.ParseMode = "Markdown"
		return sendWithRetry(bot, msg)
	}

	timezone, err := strconv.Atoi(data)
	if err != nil {
		return errors.Wrap(err, "can't parse timezone")
	}
	err = RuntimeTable.UpdateRuntimeTimezone(runtime.Uid, timezone)
	if err != nil {
		return errors.Wrap(err, "can't update timezone")
	}
	log.Info("user timezone updated", zap.Int("timezone", timezone), zap.Object("user", log.WrapUser(user)))
	Cache.DeleteRuntime(int64(user.ID))


	texts := getI18nText(runtime.Language, user, NewI18nKey(selectTimezoneSuccessfullyTextKey, timezoneToText(timezone, runtime.Language)))
	msg := botapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, texts[0])
	msg.ParseMode = "Markdown"
	return sendWithRetry(bot, msg)
}

func AddAccount(update *botapi.Update) error {
	callback := update.CallbackQuery
	user := callback.From
	runtime, err := fetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	// check account number limitation
	privilage, err := UserTable.GetUser(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't not get user from db")
	}
	if privilage.NumAccount == privilage.MaxAccount {
		texts := getI18nText(runtime.Language, user, NewI18nKey(accountReachLimitTextKey, privilage.NumAccount, privilage.MaxAccount))
		msg := botapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, texts[0])
		msg.ParseMode = "Markdown"
		return sendWithRetry(bot, msg)
	}

	proofKey, err := nintendo.NewProofKey()
	if err != nil {
		return errors.Wrap(err, "can't generate proof key")
	}
	err = Cache.SetProofKey(int64(user.ID), proofKey)
	if err != nil {
		return errors.Wrap(err, "can't save proof key")
	}
	link, err := nintendo.NewLoginLink(proofKey)
	if err != nil {
		return errors.Wrap(err, "can't generate login link")
	}
	texts := getI18nText(runtime.Language, user,
		NewI18nKey(loginLinkGuideTextKey),
		NewI18nKey(loginLinkTextKey))
	linkText := texts[1]
	msg := botapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, texts[0])
	markup := botapi.NewInlineKeyboardMarkup(botapi.NewInlineKeyboardRow(botapi.NewInlineKeyboardButtonURL(linkText, link)))
	msg.ReplyMarkup = &markup
	msg.ParseMode = "Markdown"
	return sendWithRetry(bot, msg)
}

func InputRedirectLink(update *botapi.Update) (err error) {
	text := update.Message.Text
	user := update.Message.From
	runtime, err := fetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime runtime")
	}

	code, err := nintendo.GetSessionTokenCode(text)
	if err != nil {
		// todo: invalid operation count ++
		return errors.Wrap(err, "invalid redirect link")
	}
	proofKey, err := Cache.GetProofKey(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't get proof key")
	}
	if proofKey == nil {
		// todo: invalid operation count ++
		texts := getI18nText(runtime.Language, user, NewI18nKey(expiredProofKeyTextKey))
		msg := botapi.NewMessage(update.Message.Chat.ID, texts[0])
		msg.ParseMode = "Markdown"
		return sendWithRetry(bot, msg)
	}

	texts := getI18nText(runtime.Language, user, NewI18nKey(addingAccountTextKey))
	msg := botapi.NewMessage(update.Message.Chat.ID, texts[0])
	msg.ParseMode = "Markdown"
	respMsg, err := sendWithRetryAndResponse(bot, msg)
	if err != nil {
		return err
	}

	sendFailureMessage := true
	defer func() {
		if err != nil && sendFailureMessage {
			texts := getI18nText(runtime.Language, user, NewI18nKey(addAccountUnsuccessfullyTextKey))
			msg := botapi.NewEditMessageText(update.Message.Chat.ID, respMsg.MessageID, texts[0])
			msg.ParseMode = "Markdown"
			err = sendWithRetry(bot, msg)
			if err != nil {
			}
		}
	}()

	var sessionTokes, cookie, accountName, nsName string
	err = retry(func() error {
		sessionTokes, err = nintendo.GetSessionToken(code, proofKey, runtime.Language)
		if err != nil {
			return errors.Wrap(err, "can't get session token")
		}
		cookie, accountName, nsName, err = nintendo.GetCookiesAndNames(sessionTokes, runtime.Language)
		if err != nil {
			return errors.Wrap(err, "can't get cookies")
		}
		return nil
	}, retryTimes)
	if err != nil {
		return err
	}

	// add account to db
	account := &db.Account{
		Uid:          runtime.Uid,
		SessionToken: sessionTokes,
		Tag:          accountName + ":" + nsName,
	}
	existed, err := AccountTable.IsAccountExisted(account.Uid, account.Tag)
	if err != nil {
		return errors.Wrap(err, "can't check whether account is existed")
	}
	if existed {
		texts = getI18nText(runtime.Language, user, NewI18nKey(addAccountExistedTextKey, account.Tag))
		editMsg := botapi.NewEditMessageText(update.Message.Chat.ID, respMsg.MessageID, texts[0])
		editMsg.ParseMode = "Markdown"
		return sendWithRetry(bot, editMsg)
	}
	err = Transactions.AddNewAccount(account)
	if err != nil {
		return errors.Wrap(err, "can't add new account to db")
	}
	sendFailureMessage = false

	// new user
	if runtime.SessionToken == "" {
		runtime.SessionToken = sessionTokes
		runtime.IKSM = cookie
		err = RuntimeTable.UpdateRuntimeAccount(runtime)
		if err != nil {
			log.Warn("can't update runtime runtime")
			return errors.Wrap(err, "can't update runtime runtime")
		}
		Cache.DeleteRuntime(int64(user.ID))
	}

	// notify job scheduler
	tryStartJobScheduler()

	texts = getI18nText(runtime.Language, user, NewI18nKey(addAccountSuccessfullyTextKey, account.Tag))
	editMsg := botapi.NewEditMessageText(update.Message.Chat.ID, respMsg.MessageID, texts[0])
	editMsg.ParseMode = "Markdown"
	return sendWithRetry(bot, editMsg)
}

func register(user *botapi.User) error {
	if user == nil {
		return errors.Errorf("user is nil")
	}
	newUser := &db.User{
		Uid:          int64(user.ID),
		UserName:     user.UserName,
		IsBlock:      false,
		MaxAccount:   userMaxAccount,
		NumAccount:   0,
		IsAdmin:      int64(user.ID) == defaultAdmin,
		AllowPolling: userAllowPolling,
	}
	runtime := &db.Runtime{
		Uid:          int64(user.ID),
		SessionToken: "",
		IKSM:         "0000000000000000000000000000000000000000",
		Language:     "en",
		Timezone:     480, // default UTC+8
	}
	err := Transactions.InsertUserAndRuntime(newUser, runtime)
	if err != nil {
		return errors.Wrap(err, "can't insert user and runtime to db")
	}
	if newUser.IsAdmin {
		admins.Add(newUser.Uid)
	}
	// set cache
	err = Cache.SetRuntime(runtime)
	if err != nil {
		log.Warn("can't set Runtime to cache", zap.Object("runtime", runtime), zap.Error(err))
	}
	return nil
}
