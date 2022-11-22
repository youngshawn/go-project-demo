package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DynamicMySQLCredsConfig struct {
	Username string
	Password string
}

func GetDynamicMySQLCredsConfig() *DynamicMySQLCredsConfig {
	ConfigLocker.RLock()
	defer ConfigLocker.RUnlock()

	return &DynamicMySQLCredsConfig{
		Username: Config.Database.MySQL.Username,
		Password: Config.Database.MySQL.Password,
	}
}

var (
	mysqlUsername string
	mysqlPassword string
	mysqlAddress  string
	mysqlDBname   string
	mysqlOptions  string
)

func connectMySQL() *gorm.DB {

	// get configurations
	mysqlUsername = Config.Database.MySQL.Username
	mysqlPassword = Config.Database.MySQL.Password
	mysqlAddress = Config.Database.MySQL.Address
	mysqlDBname = Config.Database.MySQL.DBname
	mysqlOptions = Config.Database.MySQL.Options

	msyql_dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", mysqlUsername,
		mysqlPassword, mysqlAddress, mysqlDBname, mysqlOptions)

	db, err := gorm.Open(mysql.Open(msyql_dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
