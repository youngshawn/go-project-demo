package config

import (
	"log"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

const db_type = "sqlite"

func init() {
	// connect to database
	if db_type == "mysql" {
		db = connect_mysql()
	} else if db_type == "sqlite" {
		db = connect_sqlite()
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	// setup connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func GetDB() *gorm.DB {
	return db
}
