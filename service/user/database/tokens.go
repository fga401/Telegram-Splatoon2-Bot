package database

import (
	"telegram-splatoon2-bot/common/enum"
	"telegram-splatoon2-bot/driver/database"
)

var tokenEnum = enum.Assign(&tokens{}).(*tokens)

type tokens struct {
	Permission userTokens
	Status     statusTokens
	Account    accountTokens
}

type statusTokens struct {
	Insert                    database.Token
	SelectByUid               database.Token
	UpdateLanguage            database.Token
	UpdateTimezone            database.Token
	UpdateIKSM                database.Token
	UpdateSessionTokenAndIKSM database.Token
}

type userTokens struct {
	Insert      database.Token
	Count       database.Token
	Admins      database.Token
	SelectByUid database.Token
}

type accountTokens struct {
	Insert            database.Token
	Delete            database.Token
	SelectByUid       database.Token
	SelectByUidAndTag database.Token
}
