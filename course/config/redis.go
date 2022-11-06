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

	// get configurations
	EnableLocalCache = Config.Cache.EnableLocalCache
	EnableNullResultCache = Config.Cache.EnableNullResultCache
	CacheTTL = time.Second * time.Duration(Config.Cache.CacheTTL)

	enable_redis := Config.Cache.EnableRedis

	if enable_redis == true {
		rdb := connect_redis()
		redisCache = cache.New(&cache.Options{
			Redis:      rdb,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	} else {
		redisCache = cache.New(&cache.Options{
			Redis:      nil,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	}
}

func connect_redis() *redis.Client {

	// get configurations
	redis_address := Config.Cache.Redis.Address
	redis_password := Config.Cache.Redis.Password
	redis_db := Config.Cache.Redis.DBindex
	redis_pool_size := Config.Cache.Redis.Pool.PoolSize
	redis_max_idle_conns := Config.Cache.Redis.Pool.MaxIdleConns
	redis_conn_max_idle_time := Config.Cache.Redis.Pool.ConnMaxIdleTime
	redis_conn_max_life_time := Config.Cache.Redis.Pool.ConnMaxLifetime
	redis_max_retries := Config.Cache.Redis.MaxRetries
	redis_pool_timeout := Config.Cache.Redis.PoolTimeout
	redis_dial_timeout := Config.Cache.Redis.DialTimeout
	redis_read_timeout := Config.Cache.Redis.ReadTimeout
	redis_write_timeout := Config.Cache.Redis.WriteTimeout

	//if enable_redis == false {
	//	return nil
	//}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_address,
		Password: redis_password,
		DB:       int(redis_db),

		PoolSize:              int(redis_pool_size),
		MaxIdleConns:          int(redis_max_idle_conns),
		ConnMaxIdleTime:       time.Second * time.Duration(redis_conn_max_idle_time),
		ConnMaxLifetime:       time.Second * time.Duration(redis_conn_max_life_time),
		MaxRetries:            int(redis_max_retries),
		PoolTimeout:           time.Second * time.Duration(redis_pool_timeout),
		DialTimeout:           time.Second * time.Duration(redis_dial_timeout),
		ReadTimeout:           time.Second * time.Duration(redis_read_timeout),
		WriteTimeout:          time.Second * time.Duration(redis_write_timeout),
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
