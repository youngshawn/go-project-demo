package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	sqlite_dbname = "course.db"
)

func connect_sqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(sqlite_dbname), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
