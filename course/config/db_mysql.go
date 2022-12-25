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

func initMySQL() *gorm.DB {
	// get configurations
	mysqlUsername = Config.Database.MySQL.Username
	mysqlPassword = Config.Database.MySQL.Password
	mysqlAddress = Config.Database.MySQL.Address
	mysqlDBname = Config.Database.MySQL.DBname
	mysqlOptions = Config.Database.MySQL.Options

	// connect mysql
	d, err := connectMySQL()
	if err != nil {
		log.Fatal("MySQL connect failed:", err)
	}

	// setup database
	if err := setupDatabase(d); err != nil {
		log.Fatal("MySQL setup failed:", err)
	}

	return d
}

func connectMySQL() (*gorm.DB, error) {
	msyql_dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", mysqlUsername,
		mysqlPassword, mysqlAddress, mysqlDBname, mysqlOptions)

	return gorm.Open(mysql.Open(msyql_dsn), &gorm.Config{})
}

func DynamicMySQLConfigReload() {

	// get dynamic mysql creds
	creds := GetDynamicMySQLCredsConfig()
	if creds.Username == mysqlUsername && creds.Password == mysqlPassword {
		log.Println("MySQL creds not changed, skip mysql reload.")
		return
	}

	DBLocker.Lock()
	defer DBLocker.Unlock()

	mysqlUsername = creds.Username
	mysqlPassword = creds.Password

	// re-init mysql
	d, err := connectMySQL()
	if err != nil {
		log.Println("MySQL connect failed:", err)
		return
	}
	if err := setupDatabase(d); err != nil {
		log.Println("MySQL setup failed:", err)
		return
	}

	db = d
}
