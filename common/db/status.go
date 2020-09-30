package db

import (
	"github.com/pkg/errors"
)

type UserStatusTable struct {
	Table
}

const (
	userStatusNamedStmtOrderInsert = iota
	userStatusNamedStmtOrderGetByUid
)

var userStatusNamedStmts = map[int]string{
	userStatusNamedStmtOrderInsert:   "INSERT INTO status (uid, iksm, is_block, max_account, n_account, is_admin, allow_polling) VALUES (:uid, :iksm, :is_block, :max_account, :n_account, :is_admin, :allow_polling);",
	userStatusNamedStmtOrderGetByUid: "SELECT uid, iksm, is_block, max_account, n_account, is_admin, allow_polling FROM status WHERE uid=:uid;",
}

var userStatusStmts = map[int]string{
}

func (impl *UserStatusTable) InsertStatus(status UserStatus) error {
	stmt := impl.namedStmts[userStatusNamedStmtOrderInsert]
	_, err := stmt.Exec(status)
	if err != nil {
		err = errors.Wrap(err, "can't execute statement insert")
	}
	return err
}
