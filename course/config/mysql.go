package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	mysql_username = "user"
	mysql_password = "pass"
	mysql_address  = "127.0.0.1:3306"
	mysql_dbname   = "dbname"
	mysql_options  = "charset=utf8mb4&parseTime=True&loc=Local"
)

var Db *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", mysql_username,
		mysql_password, mysql_address, mysql_dbname, mysql_options)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
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
