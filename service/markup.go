package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/botutil"
	log "telegram-splatoon2-bot/logger"
)

var supportLanguage []string
var markups = make(map[string]map[MarkupName]botapi.InlineKeyboardMarkup)

type MarkupName int

const (
	// settings
	settingsTextKey = "What do you want to change?"
	// language
	selectLanguageTextKey             = "Please select your preferred language:"
	selectLanguageSuccessfullyTextKey = "Change language successfully! Your language is *%s* now. Use /settings to change other settings."
	// timezone
	selectTimezoneTextKey             = "Please select your timezone:"
	selectTimezoneSuccessfullyTextKey = "Change timezone successfully! Your timezone is *%s* now. Use /settings to change other settings."
	// account
	loginLinkGuideTextKey           = "Please open the following link and:\n*1.* Login;\n*2.* Right click / Long press the *<Select this account>* button;\n*3.* Copy the link address;\n*4.* Paste and send to this bot."
	expiredProofKeyTextKey          = "Your link is expired. Please use /settings to retry."
	addingAccountTextKey            = "Fetching account from Nintendo server..."
	addAccountSuccessfullyTextKey   = "Account *%s* has been added. Use /settings to change other settings."
	addAccountUnsuccessfullyTextKey = "Adding new account failed. Please use /settings to retry."
	addAccountExistedTextKey        = "Sorry, your account is existed. Use /settings to add another accounts."
	//salmon
	salmonSchedulesFutureTextKey   = "#Future"
	salmonSchedulesScheduleTextKey = "*Time*: `%s ~ %s`\n"
	salmonSchedulesNextTextKey     = "#Next"
	salmonSchedulesOpenTextKey     = "#Open"
	salmonSchedulesSoonTextKey     = "#Soon"
	salmonSchedulesDetailTextKey   = "*Time*: `%s ~ %s`\n*Stage*: %s\n*Weapons*:\n- %s\n- %s\n- %s\n- %s\n"
)
const (
	settingsKeyboard MarkupName = iota
	languageKeyboard MarkupName = iota
	timezoneKeyboard MarkupName = iota
)
const (
	AccountSettingsKeyboardPrefix   = "<set_account>"
	LanguageSettingsKeyboardPrefix  = "<set_lang>"
	TimezoneSettingsKeyboardPrefix  = "<set_timezone>"
	LanguageSelectionKeyboardPrefix = "<sel_lang>"
	TimezoneSelectionKeyboardPrefix = "<sel_timezone>"
)

var markupPreparers = map[MarkupName]func(string, *message.Printer){
	timezoneKeyboard: prepareTimezoneKeyboard,
	languageKeyboard: prepareLanguageKeyboard,
	settingsKeyboard: prepareSettingsKeyboard,
}

func prepareTimezoneKeyboard(langTag string, printer *message.Printer) {
	markups[langTag][timezoneKeyboard] = botapi.NewInlineKeyboardMarkup(
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-12 (IDLW)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-720"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-11 (SST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-660"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-10 (HST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-600"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-9:30 (MIT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-570"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-9 (AKST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-540"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-8 (PST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-480"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-7 (MST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-420"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-6 (CST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-360"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-5 (EST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-300"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-4 (AST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-240"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-3:30 (NST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-210"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-3 (BRT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-180"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-2 (FNT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-120"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC-1 (CVT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "-60"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC (GNT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "0"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+1 (CET)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "60"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+2 (EET)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "120"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+3 (MSK)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "180"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+3:30 (IRST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "210"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+4 (GST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "240"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+4:30 (AFT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "270"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+5 (PKT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "300"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+5:30 (IST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "330"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+5:45 (NPT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "345"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+6 (BHT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "360"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+6:30 (MMT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "390"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+7 (ICT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "420"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+8 (CST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "480"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+9 (JST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "540"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+9:30 (ACST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "570"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+10 (AEST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "600"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+10:30 (LHST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "630"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+11 (VUT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "660"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+12 (NZST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "720"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+12:45 (CHAST)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "765"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+13 (PHOT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "780"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("UTC+14 (LINT)"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "840"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("« Back to Settings"),
				botutil.SetCallbackQueryPrefix(TimezoneSelectionKeyboardPrefix, "BACK"))),
	)
}

func prepareSettingsKeyboard(langTag string, printer *message.Printer) {
	markups[langTag][settingsKeyboard] = botapi.NewInlineKeyboardMarkup(
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Account"),
				botutil.SetCallbackQueryPrefix(AccountSettingsKeyboardPrefix, "")),
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Language"),
				botutil.SetCallbackQueryPrefix(LanguageSettingsKeyboardPrefix, "")),
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Timezone"),
				botutil.SetCallbackQueryPrefix(TimezoneSettingsKeyboardPrefix, ""))),
	)
}

func prepareLanguageKeyboard(langTag string, printer *message.Printer) {
	markups[langTag][languageKeyboard] = botapi.NewInlineKeyboardMarkup(
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("English"),
				botutil.SetCallbackQueryPrefix(LanguageSelectionKeyboardPrefix, "en-US"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Chinese(Traditional)"),
				botutil.SetCallbackQueryPrefix(LanguageSelectionKeyboardPrefix, "zh-TW"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Chinese(Simplified)"),
				botutil.SetCallbackQueryPrefix(LanguageSelectionKeyboardPrefix, "zh-CN"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("Japanese"),
				botutil.SetCallbackQueryPrefix(LanguageSelectionKeyboardPrefix, "ja"))),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				printer.Sprintf("« Back to Settings"),
				botutil.SetCallbackQueryPrefix(LanguageSelectionKeyboardPrefix, "BACK"))),
	)
}

func initMarkup() {
	supportLanguage = viper.GetStringSlice("service.language")
	for _, plainTag := range supportLanguage {
		markups[plainTag] = make(map[MarkupName]botapi.InlineKeyboardMarkup)
		tag, err := language.Parse(plainTag)
		if err != nil {
			panic(errors.Wrap(err, "can't parse language tag"))
		}
		printer := message.NewPrinter(tag)
		for _, preparer := range markupPreparers {
			preparer(plainTag, printer)
		}
	}
}

func getStaticMarkup(name MarkupName, tag string) botapi.InlineKeyboardMarkup {
	return markups[tag][name]
}

func langToText(selectedLanguage string, lang string) string {
	tag, err := language.Parse(lang)
	if err != nil {
		log.Warn("parse language failed", zap.String("language", lang), zap.Error(err))
		tag = language.English
	}
	printer := message.NewPrinter(tag)
	// use the same keys
	// text will be changed according to json files in ./locales
	switch selectedLanguage {
	case "en":
		return printer.Sprintf("English")
	case "ja":
		return printer.Sprintf("Japanese")
	case "zh-CN":
		return printer.Sprintf("Chinese(Simplified)")
	case "zh-TW":
		return printer.Sprintf("Chinese(Traditional)")
	}
	return printer.Sprintf("English")
}

func timezoneToText(timezone int, lang string) string {
	tag, err := language.Parse(lang)
	if err != nil {
		log.Warn("parse language failed", zap.String("language", lang), zap.Error(err))
		tag = language.English
	}
	printer := message.NewPrinter(tag)
	// use the same keys
	// text will be changed according to json files in ./locales
	switch timezone {
	case -720:
		return printer.Sprintf("UTC-12 (IDLW)")
	case -660:
		return printer.Sprintf("UTC-11 (SST)")
	case -600:
		return printer.Sprintf("UTC-10 (HST)")
	case -570:
		return printer.Sprintf("UTC-9:30 (MIT)")
	case -540:
		return printer.Sprintf("UTC-9 (AKST)")
	case -480:
		return printer.Sprintf("UTC-8 (PST)")
	case -420:
		return printer.Sprintf("UTC-7 (MST)")
	case -360:
		return printer.Sprintf("UTC-6 (CST)")
	case -300:
		return printer.Sprintf("UTC-5 (EST)")
	case -240:
		return printer.Sprintf("UTC-4 (AST)")
	case -210:
		return printer.Sprintf("UTC-3:30 (NST)")
	case -180:
		return printer.Sprintf("UTC-3 (BRT)")
	case -120:
		return printer.Sprintf("UTC-2 (FNT)")
	case -60:
		return printer.Sprintf("UTC-1 (CVT)")
	case 0:
		return printer.Sprintf("UTC (GNT)")
	case 60:
		return printer.Sprintf("UTC+1 (CET)")
	case 120:
		return printer.Sprintf("UTC+2 (EET)")
	case 180:
		return printer.Sprintf("UTC+3 (MSK)")
	case 210:
		return printer.Sprintf("UTC+3:30 (IRST)")
	case 240:
		return printer.Sprintf("UTC+4 (GST)")
	case 270:
		return printer.Sprintf("UTC+4:30 (AFT)")
	case 300:
		return printer.Sprintf("UTC+5 (PKT)")
	case 330:
		return printer.Sprintf("UTC+5:30 (IST)")
	case 345:
		return printer.Sprintf("UTC+5:45 (NPT)")
	case 360:
		return printer.Sprintf("UTC+6 (BHT)")
	case 390:
		return printer.Sprintf("UTC+6:30 (MMT)")
	case 420:
		return printer.Sprintf("UTC+7 (ICT)")
	case 480:
		return printer.Sprintf("UTC+8 (CST)")
	case 540:
		return printer.Sprintf("UTC+9 (JST)")
	case 570:
		return printer.Sprintf("UTC+9:30 (ACST)")
	case 600:
		return printer.Sprintf("UTC+10 (AEST)")
	case 630:
		return printer.Sprintf("UTC+10:30 (LHST)")
	case 660:
		return printer.Sprintf("UTC+11 (VUT)")
	case 720:
		return printer.Sprintf("UTC+12 (NZST)")
	case 765:
		return printer.Sprintf("UTC+12:45 (CHAST)")
	case 780:
		return printer.Sprintf("UTC+13 (PHOT)")
	case 840:
		return printer.Sprintf("UTC+14 (LINT)")
	}
	return printer.Sprintf("UTC+8 (CT)")
}
