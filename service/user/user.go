package user

import (
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
	"telegram-splatoon2-bot/service/user/database"
)

type ID = database.UserID
type Permission = database.Permission
type Status = database.Status
type Account = database.Account

type Service interface {
	//Group(groupType GroupType) Group
	Admins() []ID

	Existed(uid ID) (bool, error)
	Register(uid ID, username string) error

	GetStatus(uid ID) (Status, error)
	UpdateStatusIKSM(uid ID, iksm string) error
	UpdateStatusTimezone(uid ID, timezone timezone.Timezone) error
	UpdateStatusLanguage(uid ID, language language.Language) error

	GetAccount(uid ID, tag string) (Account, error)
	NewLoginLink(uid ID) error
	AddAccount(uid ID, link string) error
	DeleteAccount(uid ID, tag string) error
	SwitchAccount(uid ID, tag string) error
	ListAccounts(uid ID) ([]Account, error)

	GetPermission(uid ID) (Permission, error)
}
