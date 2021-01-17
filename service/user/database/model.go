package database

import (
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

type UserID int64

type Permission struct {
	UserID       UserID `db:"uid"`
	UserName     string `db:"user_name"`
	IsBlock      bool   `db:"is_block"`
	MaxAccount   int32  `db:"max_account"`
	IsAdmin      bool   `db:"is_admin"`
	AllowPolling bool   `db:"allow_polling"`
}

type Account struct {
	UserID       UserID `db:"uid"`
	SessionToken string `db:"session_token"`
	Tag          string `db:"tag"`
}

type Status struct {
	UserID       UserID            `db:"uid"`
	SessionToken string            `db:"session_token"`
	IKSM         string            `db:"iksm"`
	Language     language.Language `db:"language"`
	Timezone     timezone.Timezone `db:"timezone"`
}
