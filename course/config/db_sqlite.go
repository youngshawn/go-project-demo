package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connect_sqlite() *gorm.DB {

	// get configurations
	sqlite_dbname := Config.Database.Sqlite.DBname

	db, err := gorm.Open(sqlite.Open(sqlite_dbname), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
