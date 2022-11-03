package config

import (
	"log"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

func DatabaseConnectAndSetup() {

	// get configurations
	db_type := Config.Database.Type
	maxIdleConns := Config.Database.Pool.MaxIdleConns
	maxOpenConns := Config.Database.Pool.MaxOpenConns
	connMaxIdleTime := Config.Database.Pool.ConnMaxIdleTime
	connMaxLifetime := Config.Database.Pool.ConnMaxLifetime

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
	sqlDB.SetMaxIdleConns(int(maxIdleConns))
	sqlDB.SetMaxOpenConns(int(maxOpenConns))
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(connMaxIdleTime))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(connMaxLifetime))
}

func GetDB() *gorm.DB {
	return db
}
