package database

import (
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

type Service interface {
	Admins() ([]UserID, error)

	Existed(uid UserID) (bool, error)
	Register(user Permission, status Status) error

	SelectStatus(uid UserID) (Status, error)
	UpdateStatusIKSM(uid UserID, iksm string) error
	UpdateStatusTimezone(uid UserID, timezone timezone.Timezone) error
	UpdateStatusLanguage(uid UserID, language language.Language) error

	SelectAccount(uid UserID, tag string) (Account, error)
	InsertAccount(account Account) error
	SwitchAccount(uid UserID, sessionToken string, iksm string) error
	InsertAndSwitchAccount(account Account, iksm string) error
	DeleteAccount(uid UserID, tag string) error
	DeleteAndSwitchAccount(uid UserID, tag string, sessionToken string,  iksm string) error
	SelectAccounts(uid UserID) ([]Account, error)

	GetPermission(uid UserID) (Permission, error)
}
