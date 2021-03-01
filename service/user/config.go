package user

import (
	"time"

	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

// DefaultPermission is the default permission of new user.
type DefaultPermission struct {
	// Admins are IDs of admins. Any user with this ID would be treated as admin.
	Admins       []ID
	// MaxAccount is the number of account that a user can have.
	MaxAccount   int32
	// AllowPolling identifies whether a user can use 'poll battle results'.
	AllowPolling bool
	// Timezone of user.
	Timezone     timezone.Timezone
	// Language of user.
	Language     language.Language
	// IsBlock identifies whether a user can use this bot.
	IsBlock      bool
}

// Config sets up the User Service.
type Config struct {
	// AccountsCacheExpiration is the TTL of the account in cache.
	AccountsCacheExpiration time.Duration
	// ProofKeyCacheExpiration is the TTL of the proof key in cache.
	ProofKeyCacheExpiration time.Duration
	// DefaultPermission is the default permission of new user.
	DefaultPermission DefaultPermission
}
