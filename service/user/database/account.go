package database

import (
	"github.com/pkg/errors"
	"telegram-splatoon2-bot/driver/database"
)

func init() {
	registerStatements([]database.Declaration{
		{
			Token:    tokenEnum.Account.Insert,
			Stmt:     "INSERT INTO account (uid, session_token, tag) VALUES (:uid, :session_token, :tag);",
			Named:    true,
			Prepared: false,
		},
		{
			Token:    tokenEnum.Account.Delete,
			Stmt:     "DELETE FROM account WHERE uid=? AND tag=?;",
			Named:    false,
			Prepared: false,
		},
		{
			Token:    tokenEnum.Account.SelectByUID,
			Stmt:     "SELECT * FROM account WHERE uid=?;",
			Named:    false,
			Prepared: false,
		},
		{
			Token:    tokenEnum.Account.SelectByUIDAndTag,
			Stmt:     "SELECT * FROM account WHERE uid=? AND tag=?;",
			Named:    false,
			Prepared: false,
		},
		{
			Token:    tokenEnum.Status.UpdateSessionTokenAndIKSM,
			Stmt:     "UPDATE status SET session_token=?, iksm=? WHERE uid=?;",
			Named:    false,
			Prepared: false,
		},
	})
}

func (svc *serviceImpl) SelectAccount(uid UserID, tag string) (Account, error) {
	var account Account
	err := svc.db.Get(tokenEnum.Account.SelectByUIDAndTag, &account, uid, tag)
	return account, err
}

func (svc *serviceImpl) InsertAccount(account Account) error {
	return svc.db.NamedExec(tokenEnum.Account.Insert, account)
}

func (svc *serviceImpl) SwitchAccount(uid UserID, sessionToken string, iksm string) error {
	return svc.db.Transact(func(tx database.Executable) error {
		if err := tx.Exec(tokenEnum.Status.UpdateSessionTokenAndIKSM, sessionToken, iksm, uid); err != nil {
			return errors.Wrap(err, "can't update session token and IKSM")
		}
		return nil
	})
}

func (svc *serviceImpl) InsertAndSwitchAccount(account Account, iksm string) error {
	return svc.db.Transact(func(tx database.Executable) error {
		if err := tx.NamedExec(tokenEnum.Account.Insert, account); err != nil {
			return errors.Wrap(err, "can't insert account")
		}
		if err := tx.Exec(tokenEnum.Status.UpdateSessionTokenAndIKSM, account.SessionToken, iksm, account.UserID); err != nil {
			return errors.Wrap(err, "can't update session token and IKSM")
		}
		return nil
	})
}

func (svc *serviceImpl) DeleteAccount(uid UserID, tag string) error {
	return svc.db.Exec(tokenEnum.Account.Delete, uid, tag)
}

func (svc *serviceImpl) DeleteAndSwitchAccount(uid UserID, tag string, sessionToken string, iksm string) error {
	return svc.db.Transact(func(tx database.Executable) error {
		if err := tx.Exec(tokenEnum.Account.Delete, uid, tag); err != nil {
			return errors.Wrap(err, "can't delete Account")
		}
		if err := tx.Exec(tokenEnum.Status.UpdateSessionTokenAndIKSM, sessionToken, iksm, uid); err != nil {
			return errors.Wrap(err, "can't update session token and IKSM")
		}
		return nil
	})
}

func (svc *serviceImpl) SelectAccounts(uid UserID) ([]Account, error) {
	accounts := make([]Account, 0)
	err := svc.db.Select(tokenEnum.Account.SelectByUID, &accounts, uid)
	return accounts, err
}
