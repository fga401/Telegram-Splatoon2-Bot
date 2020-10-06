package db

//type UserTableInterface interface {
//	TableInterface
//	InsertUser(status User) error
//	UpdateUser(user User) error
//	GetUser(uid int64) (User, error)
//	IsUserExisted(uid int64) (bool, error)
//}

type UserTableImpl struct {
	TableImpl
}

const (
	userNamedStmtInsert namedStmtName = iota
	userNamedStmtUpdate namedStmtName = iota
)

const (
	userStmtSelectByUid     stmtName = iota
	userStmtCount           stmtName = iota
	userStmtIncreaseAccount stmtName = iota
)

var userNamedStmts = map[namedStmtName]Declaration{
	userNamedStmtInsert: {false, "INSERT INTO user (uid, user_name, is_block, max_account, n_account, is_admin, allow_polling) VALUES (:uid, :user_name, :is_block, :max_account, :n_account, :is_admin, :allow_polling);"},
	userNamedStmtUpdate: {false, "UPDATE user SET is_block=:is_block, max_account=:max_account, n_account=:n_account, is_admin=:is_admin, allow_polling=:allow_polling WHERE uid=:uid;"},
}

var userStmts = map[stmtName]Declaration{
	userStmtSelectByUid:     {false, "SELECT * FROM user WHERE uid=?;"},
	userStmtCount:           {false, "SELECT count(uid) FROM user WHERE uid=?;"},
	userStmtIncreaseAccount: {false, "UPDATE user SET n_account=n_account+1 WHERE uid=?"},
}

func (impl *UserTableImpl) InsertUser(user *User) error {
	return impl.namedExec(userNamedStmtInsert, user)
}

func (impl *UserTableImpl) UpdateUser(user *User) error {
	return impl.namedExec(userNamedStmtUpdate, user)
}

func (impl *UserTableImpl) GetUser(uid int64) (*User, error) {
	user := &User{}
	err := impl.get(userStmtSelectByUid, user, uid)
	return user, err
}

func (impl *UserTableImpl) IsUserExisted(uid int64) (bool, error) {
	var count int
	err := impl.get(userStmtCount, &count, uid)
	return count >= 1, err
}
