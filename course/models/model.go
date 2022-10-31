package models

import (
	"errors"

	"github.com/go-redis/cache/v9"
	"github.com/youngshawn/go-project-demo/course/config"
	"gorm.io/gorm"
)

var db *gorm.DB
var redisCache *cache.Cache
var ErrorObjectNotFound = errors.New("ObjectNotFound")

func init() {
	db = config.GetDB()
	//db.AutoMigrate(&Course{})
	db.AutoMigrate(&Course{}, &Teacher{})

	redisCache = config.GetCache()
}
