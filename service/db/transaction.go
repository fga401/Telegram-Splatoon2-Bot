package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type TransactionImpl struct {
	db *sqlx.DB
}

func (impl *TransactionImpl) InsertUserAndRuntime(user *User, runtime *Runtime) (err error) {
	return transact(impl.db, func (tx *sqlx.Tx) error {
		if _, err = tx.NamedExec(userNamedStmts[userNamedStmtInsert].stmt, user); err != nil {
			return errors.Wrap(err, "can't insert user")
		}
		if _, err = tx.NamedExec(runtimeNamedStmts[runtimeNamedStmtInsert].stmt, runtime); err != nil {
			return errors.Wrap(err, "can't insert runtime")
		}
		return nil
	})
}

func (impl *TransactionImpl) AddNewAccount(account *Account) (err error) {
	return transact(impl.db, func (tx *sqlx.Tx) error {
		if _, err = tx.Exec(userStmts[userStmtIncreaseAccount].stmt, account.Uid); err != nil {
			return errors.Wrap(err, "can't update user's account number")
		}
		if _, err = tx.NamedExec(accountNamedStmts[accountNamedStmtInsert].stmt, account); err != nil {
			return errors.Wrap(err, "can't insert account")
		}
		return nil
	})
}

func (impl *TransactionImpl) DeleteAccount(uid int64, tag string) (err error) {
	return transact(impl.db, func (tx *sqlx.Tx) error {
		if _, err = tx.Exec(userStmts[userStmtDecreaseAccount].stmt, uid); err != nil {
			return errors.Wrap(err, "can't update user's account number")
		}
		if _, err = tx.Exec(accountStmts[accountStmtDeleteAccount].stmt, uid, tag); err != nil {
			return errors.Wrap(err, "can't delete account")
		}
		return nil
	})
}

func transact(db *sqlx.DB, txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return errors.Wrap(err, "can't init transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	err = txFunc(tx)
	return err
}