package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"os"
)

func main() {
	viper.SetConfigName(os.Getenv("CONFIG"))
	viper.SetConfigType("json")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config failed: %v", err))
	}
	url := viper.GetString("db.url")
	sqlSource := viper.GetString("db.sql")

	db, err := sql.Open("sqlite3", url)
	if err != nil {
		panic(fmt.Errorf("open database failed: %v", err))
	}
	defer db.Close()
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		panic(fmt.Errorf("new instance failed: %v", err))
	}
	m, err := migrate.NewWithDatabaseInstance(sqlSource, "ql", driver)
	if err != nil {
		panic(fmt.Errorf("new migration failed: %v", err))
	}
	if err := m.Up(); err != nil {
		panic(fmt.Errorf("up failed: %v", err))
	}
}
