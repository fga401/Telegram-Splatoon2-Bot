package main

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
)

func main() {
	viper.SetConfigName(os.Getenv("CONFIG"))
	viper.SetConfigType("json")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("can't read config", zap.Error(err))
	}
	url := viper.GetString("db.url")
	sqlSource := viper.GetString("db.sql")

	db, err := sql.Open("sqlite3", url)
	if err != nil {
		log.Panic("can't open database", zap.Error(err))
	}
	defer db.Close()
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Panic("can't create new database instance", zap.Error(err))
	}
	m, err := migrate.NewWithDatabaseInstance(sqlSource, "ql", driver)
	if err != nil {
		log.Panic("can't create new migration instance", zap.Error(err))
	}
	if err := m.Up(); err != nil {
		log.Panic("can't upgrade", zap.Error(err))
	}
}
