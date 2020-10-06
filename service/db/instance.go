package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type User struct {
	Uid          int64  `db:"uid"`
	UserName     string `db:"user_name"`
	IsBlock      bool   `db:"is_block"`
	MaxAccount   int    `db:"max_account"`
	NumAccount   int    `db:"n_account"`
	IsAdmin      bool   `db:"is_admin"`
	AllowPolling bool   `db:"allow_polling"`
}

type Account struct {
	Uid          int64  `db:"uid"`
	SessionToken string `db:"session_token"`
	Tag          string `db:"tag"`
}

type Runtime struct {
	Uid          int64  `db:"uid"`
	SessionToken string `db:"session_token"`
	IKSM         []byte `db:"iksm"`
	Language     string `db:"language"`
}

var (
	UserTable    *UserTableImpl
	AccountTable *AccountTableImpl
	RuntimeTable *RuntimeTableImpl
	Transactions *TransactionImpl
)

type Declaration struct {
	prepared bool
	stmt     string
}

//type TableInterface interface {
//	mustPrepare(tableName string, preparedNamedStmts map[int]Declaration, preparedStmts map[int]Declaration)
//}

type namedStmtName int
type stmtName int

type TableImpl struct {
	db                    *sqlx.DB
	preparedNamedStmts    map[namedStmtName]*sqlx.NamedStmt
	preparedStmts         map[stmtName]*sqlx.Stmt
	namedStmtsDeclaration map[namedStmtName]Declaration
	stmtsDeclaration      map[stmtName]Declaration
}

func (impl *TableImpl) mustPrepare(tableName string) {
	var err error
	impl.preparedNamedStmts = make(map[namedStmtName]*sqlx.NamedStmt)
	for k, v := range impl.namedStmtsDeclaration {
		if v.prepared {
			impl.preparedNamedStmts[k], err = impl.db.PrepareNamed(v.stmt)
			if err != nil {
				panic(errors.Wrap(err, tableName+" prepare named statement failed"))
			}
		}
	}
	impl.preparedStmts = make(map[stmtName]*sqlx.Stmt)
	for k, v := range impl.stmtsDeclaration {
		if v.prepared {
			impl.preparedStmts[k], err = impl.db.Preparex(v.stmt)
			if err != nil {
				panic(errors.Wrap(err, tableName+" prepare statement failed"))
			}
		}
	}
}

func (impl *TableImpl) namedExec(name namedStmtName, arg interface{}) error {
	var err error
	if impl.namedStmtsDeclaration[name].prepared {
		stmt := impl.preparedNamedStmts[name]
		_, err = stmt.Exec(arg)
	} else {
		_, err = impl.db.NamedExec(impl.namedStmtsDeclaration[name].stmt, arg)
	}
	if err != nil {
		err = errors.Wrap(err, "can't execute named statement")
	}
	return err
}

func (impl *TableImpl) exec(name stmtName, args ...interface{}) error {
	var err error
	if impl.stmtsDeclaration[name].prepared {
		stmt := impl.preparedStmts[name]
		_, err = stmt.Exec(args...)
	} else {
		_, err = impl.db.Exec(impl.stmtsDeclaration[name].stmt, args...)
	}
	if err != nil {
		err = errors.Wrap(err, "can't execute statement")
	}
	return err
}

func (impl *TableImpl) get(name stmtName, dest interface{}, args ...interface{}) error {
	var err error
	if impl.stmtsDeclaration[name].prepared {
		stmt := impl.preparedStmts[name]
		err = stmt.Get(dest, args...)
	} else {
		err = impl.db.Get(dest, impl.stmtsDeclaration[name].stmt, args...)
	}
	if err != nil {
		err = errors.Wrap(err, "can't execute statement")
	}
	return err
}

func InitDatabaseInstance() {
	db := sqlx.MustOpen("sqlite3", viper.GetString("db.url"))
	db.SetMaxIdleConns(viper.GetInt("db.url.maxIdleConns"))
	db.SetMaxOpenConns(viper.GetInt("db.url.maxOpenConns"))

	Transactions = &TransactionImpl{db: db}
	UserTable = &UserTableImpl{TableImpl{db: db, stmtsDeclaration: userStmts, namedStmtsDeclaration: userNamedStmts}}
	AccountTable = &AccountTableImpl{TableImpl{db: db, stmtsDeclaration: accountStmts, namedStmtsDeclaration: accountNamedStmts}}
	RuntimeTable = &RuntimeTableImpl{TableImpl{db: db, stmtsDeclaration: runtimeStmts, namedStmtsDeclaration: runtimeNamedStmts}}

	UserTable.mustPrepare("UserTable")
	AccountTable.mustPrepare("AccountTable")
	RuntimeTable.mustPrepare("RuntimeTable")
}

func (r *Runtime) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", r.Uid)
	enc.AddString("session_token", r.SessionToken)
	enc.AddByteString("iksm,", r.IKSM)
	enc.AddString("language", r.Language)
	return nil
}
