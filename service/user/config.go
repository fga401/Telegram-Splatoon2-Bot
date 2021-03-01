package user

import (
	"time"

	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

type DefaultPermission struct {
	Admins       []ID
	MaxAccount   int32
	AllowPolling bool
	Timezone     timezone.Timezone
	Language     language.Language
	IsBlock      bool
}

type Config struct {
	AccountsCacheExpiration time.Duration
	ProofKeyCacheExpiration time.Duration
	DefaultPermission       DefaultPermission
}
