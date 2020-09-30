package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type UserStatus struct {
	Uid          int64    `db:"uid"`
	IKSM         [40]byte `db:"iksm"`
	IsBlock      bool     `db:"is_block"`
	MaxAccount   int16    `db:"max_account"`
	NumAccount   int16    `db:"n_account"`
	IsAdmin      bool     `db:"is_admin"`
	AllowPolling bool     `db:"allow_polling"`
}

type Account struct {
	Uid  int64    `db:"uid"`
	IKSM [40]byte `db:"iksm"`
}

type UserName struct {
	Uid      int64  `db:"uid"`
	UserName string `db:"user_name"`
}

var (
	UserStatusTableInstance *UserStatusTable
	AccountTableInstance    *AccountTable
	UserNameTableInstance   *UserNameTable
	Client                  *sqlx.DB
)

type Table struct {
	db         *sqlx.DB
	namedStmts []*sqlx.NamedStmt
	stmts      []*sqlx.Stmt
}

func (impl *Table) MustPrepare(tableName string, namedStmts map[int]string, stmts map[int]string) {
	var err error
	impl.namedStmts = make([]*sqlx.NamedStmt, len(namedStmts))
	for k, v := range namedStmts {
		impl.namedStmts[k], err = impl.db.PrepareNamed(v)
		if err != nil {
			panic(errors.Wrap(err, tableName+"prepare named statement failed"))
		}
	}
	impl.stmts = make([]*sqlx.Stmt, len(stmts))
	for k, v := range stmts {
		impl.namedStmts[k], err = impl.db.PrepareNamed(v)
		if err != nil {
			panic(errors.Wrap(err, tableName+"prepare statement failed"))
		}
	}
}

func InitDatabaseInstance() {
	db := sqlx.MustOpen("sqlite3", viper.GetString("db.url"))
	Client = db
	UserStatusTableInstance = &UserStatusTable{Table{db: db}}
	AccountTableInstance = &AccountTable{Table{db: db}}
	UserNameTableInstance = &UserNameTable{Table{db: db}}

	UserStatusTableInstance.MustPrepare("UserStatusTable", userStatusNamedStmts, userStatusStmts)
	AccountTableInstance.MustPrepare("UserStatusTable", accountNamedStmts, accountStmts)
	UserNameTableInstance.MustPrepare("UserStatusTable", userNameNamedStmts, userNameStmts)
}
