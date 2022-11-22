package config

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	sqliteFilename string
	sqliteOptions  string
)

func initSqlite() *gorm.DB {
	// get configurations
	sqliteFilename = Config.Database.Sqlite.Filename
	sqliteOptions = Config.Database.Sqlite.Options

	// connect sqlite
	d := connectSqlite()

	// setup database
	setupDatabase(d)

	return d
}

func connectSqlite() *gorm.DB {

	// sqlite3: PRAGMA foreign_keys = ON;
	sqlite_dsn := fmt.Sprintf("file:%s?%s", sqliteFilename, sqliteOptions)

	d, err := gorm.Open(sqlite.Open(sqlite_dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return d
}

func DynamicSqliteConfigReload() {
	return
}
