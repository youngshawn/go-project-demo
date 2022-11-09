package config

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connect_sqlite() *gorm.DB {

	// get configurations
	sqlite_filename := Config.Database.Sqlite.Filename
	sqlite_options := Config.Database.Sqlite.Options

	// sqlite3: PRAGMA foreign_keys = ON;
	sqlite_dsn := fmt.Sprintf("file:%s?%s", sqlite_filename, sqlite_options)

	db, err := gorm.Open(sqlite.Open(sqlite_dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
