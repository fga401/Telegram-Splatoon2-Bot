package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/botutil"
)

var supportLanguage []string
var markups = make(map[string]map[MarkupName]botapi.InlineKeyboardMarkup)

type MarkupName int

const (
	settingsTextKey                   = "What do you want to do with your account?"
	selectLanguageTextKey             = "Please select your preferred language?"
	selectLanguageSuccessfullyTextKey = "Change language successfully!"
	loginLinkGuideTextKey             = "Please open the following link and:\n*1.* Login;\n*2.* Right click / Long press the *<Select this account>* button;\n*3.* Copy the link address;\n*4.* Paste and send to this bot."
	expiredProofKeyTextKey            = "Your link is expired. Please use /settings to retry."
	addingAccountTextKey              = "Fetching account from Nintendo server..."
	addAccountSuccessfullyTextKey     = "Account *%s* has been added."
	addAccountUnsuccessfullyTextKey   = "Adding new account failed. Please use /settings to retry."
	addAccountExistedTextKey          = "Sorry, your account is existed."
)
const (
	settingsKeyboard MarkupName = iota
	languageKeyboard MarkupName = iota
)
const (
	SettingsKeyboardPrefix = "setting"
	LanguageKeyboardPrefix = "lang"
)

var markupPreparers = map[MarkupName]func(string, *message.Printer){
	languageKeyboard: prepareLanguageKeyboard,
	settingsKeyboard: prepareSettingsKeyboard,
}

func prepareSettingsKeyboard(langTag string, printer *message.Printer) {
	markups[langTag][settingsKeyboard] = botapi.NewInlineKeyboardMarkup(
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Change Language"),
				botutil.SetCallbackQueryPrefix(SettingsKeyboardPrefix, "lang"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Add Nintendo Account"),
				botutil.SetCallbackQueryPrefix(SettingsKeyboardPrefix, "account"))),
	)
}

func prepareLanguageKeyboard(langTag string, printer *message.Printer) {
	markups[langTag][languageKeyboard] = botapi.NewInlineKeyboardMarkup(
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("English"),
				botutil.SetCallbackQueryPrefix(LanguageKeyboardPrefix, "en-US"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Chinese(Traditional)"),
				botutil.SetCallbackQueryPrefix(LanguageKeyboardPrefix, "zh-CN"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Chinese(Simplified)"),
				botutil.SetCallbackQueryPrefix(LanguageKeyboardPrefix, "zh-TW"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Japanese"),
				botutil.SetCallbackQueryPrefix(LanguageKeyboardPrefix, "ja"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("« Back to Settings"),
				botutil.SetCallbackQueryPrefix(LanguageKeyboardPrefix, "BACK"))),
	)
}

func initMarkup() {
	supportLanguage = viper.GetStringSlice("service.language")
	for _, planeTag := range supportLanguage {
		markups[planeTag] = make(map[MarkupName]botapi.InlineKeyboardMarkup)
		tag, err := language.Parse(planeTag)
		if err != nil {
			panic(errors.Wrap(err, "can't parse language tag"))
		}
		printer := message.NewPrinter(tag)
		for _, preparer := range markupPreparers {
			preparer(planeTag, printer)
		}
	}
}

func getStaticMarkup(name MarkupName, tag string) botapi.InlineKeyboardMarkup {
	return markups[tag][name]
}
