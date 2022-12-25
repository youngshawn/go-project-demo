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
	d, err := connectSqlite()
	if err != nil {
		log.Fatal("Sqlite connect failed:", err)
	}

	// setup database
	if err := setupDatabase(d); err != nil {
		log.Println("Sqlite setup failed:", err)
	}

	return d
}

func connectSqlite() (*gorm.DB, error) {

	// sqlite3: PRAGMA foreign_keys = ON;
	sqlite_dsn := fmt.Sprintf("file:%s?%s", sqliteFilename, sqliteOptions)

	return gorm.Open(sqlite.Open(sqlite_dsn), &gorm.Config{})
}

func DynamicSqliteConfigReload() {
	return
}
