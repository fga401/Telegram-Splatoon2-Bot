package router

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/common/enum"
)

// Mode use to specify the Router running mode.
type Mode enum.Enum
type modeEnum struct {
	// Polling mode means that Router would send a 'getUpdates' request to telegram to get update periodically.
	// More info: https://core.telegram.org/bots/api#callbackquery
	Polling Mode
	// WebHook mode means that whenever there is an update for the bot, telegram will send an HTTPS POST request to the specified url.
	// More info: https://core.telegram.org/bots/api#callbackquery
	WebHook Mode
}

// ModeEnum lists all available Mode.
var ModeEnum = enum.Assign(&modeEnum{}).(*modeEnum)

// PollingConfig sets up Polling mode.
type PollingConfig = botApi.UpdateConfig

// WebHookConfig sets up WebHook mode.
type WebHookConfig struct {
	URL  string
	Cert string
	Key  string
	Port string
}

// Config sets up Router.
type Config struct {
	// Mode in which Router is running.
	Mode Mode
	// MaxWorker sets the max number of goroutine to process request.
	// If MaxWorker == 0, there is no limitation.
	MaxWorker int32

	// Polling config. Router reads this config if it's in Polling mode.
	Polling PollingConfig
	// WebHook config. Router reads this config if it's in WebHook mode.
	WebHook WebHookConfig
}
