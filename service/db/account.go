package db

type Account struct {
	Uid          int64  `db:"uid"`
	SessionToken string `db:"session_token"`
	Tag          string `db:"tag"`
}

type AccountTableImpl struct {
	TableImpl
}

const (
	accountNamedStmtInsert namedStmtName = iota
	accountNamedStmtUpdate namedStmtName = iota
)

const (
	accountStmtSelectAccounts stmtName = iota
	accountStmtSelectAccount  stmtName = iota
	accountStmtDeleteAccount  stmtName = iota
	accountStmtCount          stmtName = iota
)

var accountNamedStmts = map[namedStmtName]Declaration{
	accountNamedStmtInsert: {false, "INSERT INTO account (uid, session_token, tag) VALUES (:uid, :session_token, :tag);"},
	accountNamedStmtUpdate: {false, "UPDATE account SET session_token=:session_token, tag=:tag WHERE uid=:uid;"},
}

var accountStmts = map[stmtName]Declaration{
	accountStmtSelectAccounts: {false, "SELECT * FROM account WHERE uid=?;"},
	accountStmtSelectAccount:  {false, "SELECT * FROM account WHERE uid=? AND tag=?;"},
	accountStmtDeleteAccount:  {false, "DELETE FROM account WHERE uid=? AND tag=?;"},
	accountStmtCount:          {false, "SELECT count(tag) FROM account WHERE uid=? AND tag=?;"},
}

func (impl *AccountTableImpl) InsertAccount(account *Account) error {
	return impl.namedExec(accountNamedStmtInsert, account)
}

func (impl *AccountTableImpl) UpdateAccount(account *Account) error {
	return impl.namedExec(accountNamedStmtUpdate, account)
}

func (impl *AccountTableImpl) GetAccount(uid int64, tag string) (*Account, error) {
	account := &Account{}
	err := impl.get(accountStmtSelectAccount, account, uid, tag)
	return account, err
}

func (impl *AccountTableImpl) GetAccounts(uid int64) ([]*Account, error) {
	account := make([]*Account, 0)
	err := impl.sel(accountStmtSelectAccounts, &account, uid)
	return account, err
}

func (impl *AccountTableImpl) IsAccountExisted(uid int64, tag string) (bool, error) {
	var count int
	err := impl.get(accountStmtCount, &count, uid, tag)
	return count >= 1, err
}
