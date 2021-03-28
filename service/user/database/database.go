package database

import (
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

// Service Interacts with the database and manages users.
type Service interface {
	// Admins loads all admin UserIDs.
	Admins() ([]UserID, error)

	// Existed checks whether a user is existed.
	Existed(uid UserID) (bool, error)
	// Register adds a new user to database.
	Register(user User, permission Permission, status Status) error

	// SelectStatus gets the status against the user.
	SelectStatus(uid UserID) (Status, error)
	// UpdateStatusIKSM updates the IKSM of user.
	UpdateStatusIKSM(uid UserID, iksm string) error
	// UpdateStatusTimezone updates the timezone of user.
	UpdateStatusTimezone(uid UserID, timezone timezone.Timezone) error
	// UpdateStatusLanguage updates the language of user.
	UpdateStatusLanguage(uid UserID, language language.Language) error
	// UpdateStatusLastBattle updates the lastBattle of user.
	UpdateStatusLastBattle(uid UserID, lastBattle string) error
	// UpdateStatusLastSalmon updates the lastSalmon of user.
	UpdateStatusLastSalmon(uid UserID, lastSalmon string) error

	// SelectStatus gets the account against the user.
	SelectAccount(uid UserID, tag string) (Account, error)
	// InsertAccount adds a new account to the user.
	InsertAccount(account Account) error
	// SwitchAccount switches the user account.
	SwitchAccount(uid UserID, sessionToken string, iksm string) error
	// InsertAndSwitchAccount adds a new account to the user and switches to it.
	InsertAndSwitchAccount(account Account, iksm string) error
	// DeleteAccount deletes a user account.
	DeleteAccount(uid UserID, tag string) error
	// DeleteAccount deletes a user account and switch to a new one given sessionToken and IKSM.
	DeleteAndSwitchAccount(uid UserID, tag string, sessionToken string, iksm string) error
	// SelectAccounts loads all accounts of the user.
	SelectAccounts(uid UserID) ([]Account, error)

	// GetPermission gets the permission against the user.
	GetPermission(uid UserID) (Permission, error)
}
