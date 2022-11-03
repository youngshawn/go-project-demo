package config

import (
	"context"
	"log"
	"time"

	cache "github.com/go-redis/cache/v9"
	redis "github.com/go-redis/redis/v9"
)

var redisCache *cache.Cache
var EnableNullResultCache bool
var EnableLocalCache bool
var CacheTTL time.Duration

func CacheConnectAndSetup() {

	rdb := connect_redis()
	redisCache = cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
}

func connect_redis() *redis.Client {

	// get configurations
	EnableNullResultCache = Config.Cache.EnableNullResultCache
	EnableLocalCache = Config.Cache.EnableLocalCache
	CacheTTL = time.Second * time.Duration(Config.Cache.CacheTTL)
	redis_address := Config.Cache.Redis.Address
	redis_password := Config.Cache.Redis.Password
	redis_db := Config.Cache.Redis.DBindex

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_address,
		Password: redis_password,
		DB:       int(redis_db),

		PoolSize:              100,
		PoolTimeout:           time.Second * 10,
		ConnMaxIdleTime:       time.Minute * 10,
		ConnMaxLifetime:       time.Hour,
		MaxRetries:            1,
		MaxIdleConns:          10,
		DialTimeout:           time.Second * 5,
		ReadTimeout:           time.Second * 5,
		WriteTimeout:          time.Second * 3,
		ContextTimeoutEnabled: true,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	return rdb
}

func GetCache() *cache.Cache {
	return redisCache
}
