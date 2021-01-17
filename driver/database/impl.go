package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
)

type databaseImpl struct {
	db *sqlx.DB

	preparedNamedStmts map[Token]*sqlx.NamedStmt
	preparedStmts      map[Token]*sqlx.Stmt
	stmts              map[Token]Declaration
}

func New(config Config) Database {
	db := sqlx.MustOpen(config.Driver, config.URL)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	driver := &databaseImpl{
		db:                 db,
		preparedNamedStmts: make(map[Token]*sqlx.NamedStmt),
		preparedStmts:      make(map[Token]*sqlx.Stmt),
		stmts:              make(map[Token]Declaration),
	}
	return driver
}

func (impl *databaseImpl) MustPrepare(stmts []Declaration) {
	var err error
	for _, s := range stmts {
		if _, found := impl.stmts[s.Token]; found {
			log.Panic("token already existed")
		}
		impl.stmts[s.Token] = s
		if s.Prepared {
			if s.Named {
				impl.preparedNamedStmts[s.Token], err = impl.db.PrepareNamed(s.Stmt)
				if err != nil {
					log.Panic("can't prepare named statement", zap.Error(err))
				}
			} else {
				impl.preparedStmts[s.Token], err = impl.db.Preparex(s.Stmt)
				if err != nil {
					log.Panic("can't prepare named statement", zap.Error(err))
				}
			}
		}
	}
}

func (impl *databaseImpl) NamedExec(token Token, arg interface{}) error {
	var err error
	if impl.stmts[token].Prepared {
		stmt := impl.preparedNamedStmts[token]
		_, err = stmt.Exec(arg)
	} else {
		_, err = impl.db.NamedExec(impl.stmts[token].Stmt, arg)
	}
	if err != nil {
		err = errors.Wrap(err, "can't execute named statement")
	}
	return err
}

func (impl *databaseImpl) Exec(token Token, args ...interface{}) error {
	var err error
	if impl.stmts[token].Prepared {
		stmt := impl.preparedStmts[token]
		_, err = stmt.Exec(args...)
	} else {
		_, err = impl.db.Exec(impl.stmts[token].Stmt, args...)
	}
	if err != nil {
		err = errors.Wrap(err, "can't execute statement")
	}
	return err
}

func (impl *databaseImpl) Get(token Token, dest interface{}, args ...interface{}) error {
	var err error
	if impl.stmts[token].Prepared {
		stmt := impl.preparedStmts[token]
		err = stmt.Get(dest, args...)
	} else {
		err = impl.db.Get(dest, impl.stmts[token].Stmt, args...)
	}
	if err != nil {
		err = errors.Wrap(err, "can't execute get statement")
	}
	return err
}

func (impl *databaseImpl) Select(token Token, dest interface{}, args ...interface{}) error {
	var err error
	if impl.stmts[token].Prepared {
		stmt := impl.preparedStmts[token]
		err = stmt.Select(dest, args...)
	} else {
		err = impl.db.Select(dest, impl.stmts[token].Stmt, args...)
	}
	if err != nil {
		err = errors.Wrap(err, "can't execute select statement")
	}
	return err
}

func (impl *databaseImpl) Transact(txFunc func(tx Executable) error) error {
	tx, err := impl.db.Beginx()
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
	err = txFunc(&txAdapter{
		tx:    tx,
		table: impl,
	})
	return err
}
