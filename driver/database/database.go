package database

import "telegram-splatoon2-bot/common/enum"

// Declaration of a SQL statement.
type Declaration struct {
	// Token is the unique ID of a statement.
	Token Token
	// Stmt, the SQL statement.
	Stmt string
	// Named identifies if the statement is named statement, i.e. the statement containing any named placeholder parameters which are replaced with fields from arg.
	Named bool
	// Prepared identifies whether to prepare the statement. If prepared, a connection will be kept for this statement.
	Prepared bool
}

// Token is the unique ID of a statement.
type Token enum.Enum

// Executable wraps a database or a transaction instance.
type Executable interface {
	NamedExec(token Token, arg interface{}) error
	Exec(token Token, args ...interface{}) error
	Get(token Token, dest interface{}, args ...interface{}) error
	Select(token Token, dest interface{}, args ...interface{}) error
}

// Database manages and executes all SQL statement and transaction.
type Database interface {
	Executable
	// Transact executes the transaction txFunc. If txFunc returns an error, transaction will be rollback.
	Transact(txFunc func(tx Executable) error) error
	// MustPrepare prepares all statements.
	MustPrepare(stmts []Declaration)
}
