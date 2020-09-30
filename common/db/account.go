package db

import "github.com/pkg/errors"

type AccountTable struct {
	Table
}

const (
	accountNamedStmtOrderInsert = iota
)

var accountNamedStmts = map[int]string{
	accountNamedStmtOrderInsert: "INSERT INTO account (uid, iksm) VALUES (:uid, :iksm);",
}

var accountStmts = map[int]string{
}

func (impl *UserNameTable) InsertAccount(account Account) error {
	stmt := impl.namedStmts[accountNamedStmtOrderInsert]
	_, err := stmt.Exec(account)
	if err != nil {
		err = errors.Wrap(err, "can't execute statement insert")
	}
	return err
}
