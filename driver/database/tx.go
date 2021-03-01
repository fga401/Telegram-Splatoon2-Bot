package database

import "github.com/jmoiron/sqlx"

type txAdapter struct {
	tx    *sqlx.Tx
	table *databaseImpl
}

func (w *txAdapter) NamedExec(token Token, arg interface{}) error {
	_, err := w.tx.NamedExec(w.tokenToStmt(token), arg)
	return err
}

func (w *txAdapter) Exec(token Token, args ...interface{}) error {
	_, err := w.tx.Exec(w.tokenToStmt(token), args...)
	return err
}

func (w *txAdapter) Get(token Token, dest interface{}, args ...interface{}) error {
	err := w.tx.Get(dest, w.tokenToStmt(token), args...)
	return err
}

func (w *txAdapter) Select(token Token, dest interface{}, args ...interface{}) error {
	err := w.tx.Select(dest, w.tokenToStmt(token), args...)
	return err
}

func (w *txAdapter) tokenToStmt(token Token) string {
	return w.table.stmts[token].Stmt
}
