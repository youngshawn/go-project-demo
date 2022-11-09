package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connect_mysql() *gorm.DB {

	// get configurations
	mysql_username := Config.Database.MySQL.Username
	mysql_password := Config.Database.MySQL.Password
	mysql_address := Config.Database.MySQL.Address
	mysql_dbname := Config.Database.MySQL.DBname
	mysql_options := Config.Database.MySQL.Options

	msyql_dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", mysql_username,
		mysql_password, mysql_address, mysql_dbname, mysql_options)

	db, err := gorm.Open(mysql.Open(msyql_dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
