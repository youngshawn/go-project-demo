package models

import (
	"errors"

	"github.com/go-redis/cache/v9"
	vault "github.com/hashicorp/vault/api"
	"github.com/youngshawn/go-project-demo/course/config"
	"gorm.io/gorm"
)

var db *gorm.DB
var Cache *cache.Cache
var Vault *vault.Client
var ErrorObjectNotFound = errors.New("ObjectNotFound")

func ModelInit() {
	db = config.GetDB()
	//db.AutoMigrate(&Course{})
	db.AutoMigrate(&Course{}, &Teacher{})

	Cache = config.GetCache()
	Vault = config.GetVaultClient()
}

func SetDB(d *gorm.DB) {
	db = d
}
