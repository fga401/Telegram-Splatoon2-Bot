package router

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-splatoon2-bot/common/enum"
)

type Mode enum.Enum
type modeEnum struct {
	Polling Mode
	WebHook Mode
}

var ModeEnum = enum.Assign(&modeEnum{}).(*modeEnum)

type PollingConfig = botApi.UpdateConfig

type WebHookConfig struct {
	Url  string
	Cert string
	Key  string
	Port string
}

type Config struct {
	Mode      Mode
	MaxWorker int32

	Polling PollingConfig
	WebHook WebHookConfig
}
