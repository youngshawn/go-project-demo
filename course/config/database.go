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

	// connect and setup
	databaseConnectAndSetup()
}

func DynamicMySQLCredsConfigHandler() {
	if dbType != "mysql" {
		return
	}
	// get dynamic mysql creds
	creds := GetDynamicMySQLCredsConfig()
	if creds.Username == mysqlUsername && creds.Password == mysqlPassword {
		return
	}

	DBLocker.Lock()
	defer DBLocker.Unlock()

	mysqlUsername = creds.Username
	mysqlPassword = creds.Password

	databaseConnectAndSetup()
}

func databaseConnectAndSetup() {
	// connect to database
	if dbType == "mysql" {
		db = connectMySQL()
	} else if dbType == "sqlite" {
		db = connectSqlite()
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
