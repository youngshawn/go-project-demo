package config

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectSqlite() *gorm.DB {

	// get configurations
	sqliteFilename := Config.Database.Sqlite.Filename
	sqliteOptions := Config.Database.Sqlite.Options

	// sqlite3: PRAGMA foreign_keys = ON;
	sqlite_dsn := fmt.Sprintf("file:%s?%s", sqliteFilename, sqliteOptions)

	db, err := gorm.Open(sqlite.Open(sqlite_dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
