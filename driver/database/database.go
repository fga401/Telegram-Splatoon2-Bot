package database

import "telegram-splatoon2-bot/common/enum"

type Declaration struct {
	Token    Token
	Stmt     string
	Named    bool
	Prepared bool
}
type Token enum.Enum
type Executable interface {
	NamedExec(token Token, arg interface{}) error
	Exec(token Token, args ...interface{}) error
	Get(token Token, dest interface{}, args ...interface{}) error
	Select(token Token, dest interface{}, args ...interface{}) error
}
type Database interface {
	Executable
	Transact(txFunc func(tx Executable) error) error
	MustPrepare(stmts []Declaration)
}
