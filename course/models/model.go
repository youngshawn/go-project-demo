package models

import (
	"github.com/youngshawn/go-project-demo/course/config"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = config.GetDB()
	//db.AutoMigrate(&Course{})
	db.AutoMigrate(&Course{}, &Teacher{})
}
