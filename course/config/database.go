package config

import (
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB
var DBLocker sync.Mutex

var (
	dbType          string
	maxIdleConns    uint
	maxOpenConns    uint
	connMaxIdleTime uint
	connMaxLifetime uint
)

func DatabaseInit() {
	// get configurations
	dbType = Config.Database.Type
	maxIdleConns = Config.Database.Pool.MaxIdleConns
	maxOpenConns = Config.Database.Pool.MaxOpenConns
	connMaxIdleTime = Config.Database.Pool.ConnMaxIdleTime
	connMaxLifetime = Config.Database.Pool.ConnMaxLifetime

	// init database
	if dbType == "mysql" {
		db = initMySQL()
	} else if dbType == "sqlite" {
		db = initSqlite()
	} else {
		log.Fatalf("database.type (%s) not supported.\n", dbType)
	}
}

func setupDatabase(db *gorm.DB) error {
	// get sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// setup connection pool
	sqlDB.SetMaxIdleConns(int(maxIdleConns))
	sqlDB.SetMaxOpenConns(int(maxOpenConns))
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(connMaxIdleTime))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(connMaxLifetime))
	return nil
}

func DynamicDatabaseConfigReload() {
	if dbType == "mysql" {
		DynamicMySQLConfigReload()
	}
	if dbType == "sqlite" {
		DynamicSqliteConfigReload()
	}
}

func GetDB() *gorm.DB {
	return db
}
