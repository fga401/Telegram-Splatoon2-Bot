package db

import "github.com/pkg/errors"

type UserNameTable struct {
	Table
}

const (
	userNameNamedStmtOrderInsert = iota
)

var userNameNamedStmts = map[int]string{
	userNameNamedStmtOrderInsert: "INSERT INTO username (uid, user_name) VALUES (:uid, :user_name);",
}

var userNameStmts = map[int]string{
}

func (impl *UserNameTable) InsertUser(user UserName) error {
	stmt := impl.namedStmts[userNameNamedStmtOrderInsert]
	_, err := stmt.Exec(user)
	if err != nil {
		err = errors.Wrap(err, "can't execute statement insert")
	}
	return err
}
