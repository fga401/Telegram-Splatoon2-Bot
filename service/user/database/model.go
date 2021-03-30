package database

import (
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

// UserID is ID of user.
type UserID int64

// Permission database structure storing user permission.
type Permission struct {
	UserID       UserID `db:"uid"`
	IsBlock      bool   `db:"is_block"`
	MaxAccount   int32  `db:"max_account"`
	IsAdmin      bool   `db:"is_admin"`
	AllowPolling bool   `db:"allow_polling"`
}

// Account database structure storing user accounts.
type Account struct {
	UserID       UserID `db:"uid"`
	SessionToken string `db:"session_token"`
	Tag          string `db:"tag"`
}

// Status database structure storing user status and preference.
type Status struct {
	UserID       UserID            `db:"uid"`
	SessionToken string            `db:"session_token"`
	IKSM         string            `db:"iksm"`
	LastBattle   string            `db:"last_battle"`
	LastSalmon   string            `db:"last_salmon"`
	Language     language.Language `db:"language"`
	Timezone     timezone.Timezone `db:"timezone"`
}

// User database structure storing userID and userName.
type User struct {
	UserID   UserID `db:"uid"`
	UserName string `db:"user_name"`
}
