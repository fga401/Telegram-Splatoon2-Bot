package database

import "telegram-splatoon2-bot/driver/database"

func init() {
	registerStatements([]database.Declaration{
		{
			Token:    tokenEnum.Permission.SelectByUID,
			Stmt:     "SELECT * FROM permission WHERE uid=?;",
			Named:    false,
			Prepared: false,
		},
	})
}

func (svc *serviceImpl) GetPermission(uid UserID) (Permission, error) {
	var user Permission
	err := svc.db.Get(tokenEnum.Permission.SelectByUID, &user, uid)
	return user, err
}
