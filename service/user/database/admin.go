package database

import (
	"telegram-splatoon2-bot/driver/database"
)

func init() {
	registerStatements([]database.Declaration{
		database.Declaration{
			Token:    tokenEnum.Permission.Admins,
			Stmt:     "SELECT uid FROM permission WHERE is_admin=true;",
			Named:    false,
			Prepared: false,
		},
	})
}

func (svc *serviceImpl) Admins() ([]UserID, error) {
	var ret []UserID
	err := svc.db.Select(tokenEnum.Permission.Admins, &ret)
	return ret, err
}
