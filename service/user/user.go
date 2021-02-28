package user

import (
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
	"telegram-splatoon2-bot/service/user/database"
)

// ID of user
type ID = database.UserID
// Permission stores user permission.
type Permission = database.Permission
// Status stores user status and preference.
type Status = database.Status
// Account stores user accounts.
type Account = database.Account

// Service manages all transactions about user.
type Service interface {
	// Admins loads all admin UserIDs.
	Admins() []ID

	// Existed checks whether a user is existed.
	Existed(uid ID) (bool, error)
	// Register adds a new user to database.
	Register(uid ID, username string) error

	// GetStatus gets the status against the user.
	GetStatus(uid ID) (Status, error)
	// UpdateStatusIKSM updates the IKSM of user and return the new status.
	UpdateStatusIKSM(uid ID) (Status, error)
	// UpdateStatusTimezone updates the timezone of user and return the new status.
	UpdateStatusTimezone(uid ID, timezone timezone.Timezone) (Status, error)
	// UpdateStatusLanguage updates the language of user and return the new status.
	UpdateStatusLanguage(uid ID, language language.Language) (Status, error)

	// GetAccount gets the account against the user.
	GetAccount(uid ID, tag string) (Account, error)
	// NewLoginLink generates a login link to the user.
	NewLoginLink(uid ID) (string, error)
	// AddAccount add an account to the user by the input link.
	// If user has no account, it will switch the new account.
	AddAccount(uid ID, link string) (Account, error)
	// DeleteAccount deletes a user account.
	// If delete the current account, it will switch the first account in list.
	DeleteAccount(uid ID, tag string) error
	// SwitchAccount switches the user account.
	SwitchAccount(uid ID, tag string) error
	// ListAccounts loads all accounts of the user.
	ListAccounts(uid ID) ([]Account, error)

	// GetPermission gets the permission against the user.
	GetPermission(uid ID) (Permission, error)
}
