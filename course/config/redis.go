package config

import (
	"context"
	"time"

	cache "github.com/go-redis/cache/v9"
	redis "github.com/go-redis/redis/v9"
)

const (
	redis_address  = "127.0.0.1:6379"
	redis_password = ""
	redis_db       = 0
)

var redisCache *cache.Cache

// Dynamic Config
var EnableNullResultCache bool = false
var EnableLocalCache bool = false
var CacheTTL time.Duration = time.Hour

func init() {
	rdb := connect_redis()
	redisCache = cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
}

func connect_redis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_address,
		Password: redis_password,
		DB:       redis_db,

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
		panic(err)
	}

	return rdb
}

func GetCache() *cache.Cache {
	return redisCache
}
