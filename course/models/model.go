package models

import (
	"context"
	"errors"
	"time"

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

func Healthcheck() (status string, details map[string]string) {

	details = make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// DB check
	details["database"] = "unhealthy"
	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.PingContext(ctx); err == nil {
			details["database"] = "healthy"
		}
	}

	// Redis check
	details["redis"] = "unhealthy"
	testKey := "__redis_status__"
	testValue := "healthy"
	var test string
	if err := Cache.Delete(ctx, testKey); err == nil {
		if err := Cache.Set(&cache.Item{
			Ctx:            ctx,
			Key:            testKey,
			Value:          testValue,
			TTL:            time.Hour,
			SkipLocalCache: true,
		}); err == nil {
			if err := Cache.GetSkippingLocalCache(ctx, testKey, &test); err == nil {
				if test == testValue {
					details["redis"] = "healthy"
				}
			}
		}
	}

	// Vault Check
	details["vault"] = "unhealthy"
	testPlainIDcard := "plain id card for test"
	testPlainPhone := "plain phone for test"
	teacher := &Teacher{
		PlainIDcard: testPlainIDcard,
		PlainPhone:  testPlainPhone,
	}
	if err := teacher.Encrypt(); err == nil {
		teacher.PlainIDcard = ""
		teacher.PlainPhone = ""
		if err := teacher.Decrypt(); err == nil {
			if teacher.PlainIDcard == testPlainIDcard && teacher.PlainPhone == testPlainPhone {
				details["vault"] = "healthy"
			}
		}
	}

	// summarize
	if details["database"] == "unhealthy" || details["vault"] == "unhealthy" {
		status = "unhealthy"
	} else {
		status = "healthy"
	}

	return
}
