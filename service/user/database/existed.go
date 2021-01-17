package database

import (
	"telegram-splatoon2-bot/driver/database"
)

func init() {
	registerStatements([]database.Declaration{
		database.Declaration{
			Token:    tokenEnum.Permission.Count,
			Stmt:     "SELECT count(uid) FROM permission WHERE uid=?;",
			Named:    false,
			Prepared: false,
		},
	})
}

func (svc *serviceImpl) Existed(uid UserID) (bool, error) {
	var count int
	err := svc.db.Get(tokenEnum.Permission.Count, &count, uid)
	return count >= 0, err
}
