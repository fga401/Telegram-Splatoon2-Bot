package database

import (
	"telegram-splatoon2-bot/driver/database"
)

type serviceImpl struct {
	db database.Database
}

func New(db database.Database) Service {
	svc:= &serviceImpl{
		db: db,
	}
	svc.db.MustPrepare(statement)
	return svc
}

var statement = make([]database.Declaration, 0)

func registerStatements(stmts []database.Declaration) {
	statement = append(statement, stmts...)
}
