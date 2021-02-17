package database

import (
	"github.com/pkg/errors"
	"telegram-splatoon2-bot/driver/database"
)

func init() {
	registerStatements([]database.Declaration{
		{
			Token:    tokenEnum.Permission.Insert,
			Stmt:     "INSERT INTO permission (uid, is_block, max_account, is_admin, allow_polling) VALUES (:uid, :is_block, :max_account, :is_admin, :allow_polling);",
			Named:    true,
			Prepared: false,
		},
		{
			Token:    tokenEnum.Status.Insert,
			Stmt:     "INSERT INTO status (uid, session_token, iksm, language, timezone) VALUES (:uid, :session_token, :iksm, :language, :timezone);",
			Named:    true,
			Prepared: false,
		},
	})
}

func (svc *serviceImpl) Register(user Permission, status Status) error {
	return svc.db.Transact(func(tx database.Executable) error {
		if err := tx.NamedExec(tokenEnum.Permission.Insert, user); err != nil {
			return errors.Wrap(err, "can't insert Permission")
		}
		if err := tx.NamedExec(tokenEnum.Status.Insert, status); err != nil {
			return errors.Wrap(err, "can't insert Status")
		}
		return nil
	})
}
