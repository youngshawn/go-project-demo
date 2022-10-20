package models

import (
	"github.com/youngshawn/go-project-demo/course/config"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = config.Db
	db.AutoMigrate(&Course{}, &Teacher{})
}
