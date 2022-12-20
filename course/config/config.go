package config

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Config config
var ConfigLocker sync.RWMutex
var ViperLocker sync.Mutex

var (
	Version   string
	GitCommit string
	GoVersion string
	OsArch    string
)

type config struct {
	Listen   string
	Database struct {
		Type   string
		Sqlite struct {
			Filename string
			Options  string
		}
		MySQL struct {
			Username string
			Password string
			Address  string
			DBname   string
			Options  string
		}
		Pool struct {
			MaxIdleConns    uint `mapstructure:"max-idle-conns"`
			MaxOpenConns    uint `mapstructure:"max-open-conns"`
			ConnMaxIdleTime uint `mapstructure:"conn-max-idle-time"`
			ConnMaxLifetime uint `mapstructure:"conn-max-life-time"`
		}
	}
	Cache struct {
		EnableRedis           bool `mapstructure:"enable-redis"`
		EnableLocalCache      bool `mapstructure:"enable-local-cache"`       //Dynamic Config
		EnableNullResultCache bool `mapstructure:"enable-null-result-cache"` //Dynamic Config
		CacheTTL              uint `mapstructure:"cache-ttl"`                //Dynamic Config
		Redis                 struct {
			Address  string
			Password string
			DBindex  uint `mapstructure:"db"`
			Pool     struct {
				PoolSize        uint `mapstructure:"pool-size"`
				MaxIdleConns    uint `mapstructure:"max-idle-conns"`
				ConnMaxIdleTime uint `mapstructure:"conn-max-idle-time"`
				ConnMaxLifetime uint `mapstructure:"conn-max-life-time"`
			}
			MaxRetries   uint `mapstructure:"max-retries"`
			PoolTimeout  uint `mapstructure:"pool-timeout"`
			DialTimeout  uint `mapstructure:"dial-timeout"`
			ReadTimeout  uint `mapstructure:"read-timeout"`
			WriteTimeout uint `mapstructure:"write-timeout"`
		}
	}
	RemoteConfig struct {
		Enable   bool
		Provider string
		Endpoint string
		Path     string
		Format   string
	}
	Vault struct {
		Address string
		Auth    struct {
			RoleIdFilePath   string `mapstructure:"roleid-file-path"`
			SecretIdFilePath string `mapstructure:"secretid-file-path"`
			Wrapped          bool
		}
		Transit struct {
			Key string
		}
	}
}

func init() {
	// viper set defaults
	viper.SetDefault("listen", ":8080")
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.sqlite.filename", "course.db")
	viper.SetDefault("database.sqlite.options", "_foreign_keys=on")
	viper.SetDefault("database.mysql.username", "username")
	viper.SetDefault("database.mysql.password", "password")
	viper.SetDefault("database.mysql.address", "127.0.0.1:3306")
	viper.SetDefault("database.mysql.dbname", "course")
	viper.SetDefault("database.mysql.options", "charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("database.pool.max-idle-conns", 10)
	viper.SetDefault("database.pool.max-open-conns", 100)
	viper.SetDefault("database.pool.conn-max-idle-time", 300)
	viper.SetDefault("database.pool.conn-max-life-time", 3600)
	viper.SetDefault("cache.enable-redis", true)
	viper.SetDefault("cache.enable-local-cache", true)
	viper.SetDefault("cache.enable-null-result-cache", true)
	viper.SetDefault("cache.cache-ttl", 3600)
	viper.SetDefault("cache.redis.address", "127.0.0.1:6379")
	viper.SetDefault("cache.redis.password", "")
	viper.SetDefault("cache.redis.db", 0)
	viper.SetDefault("cache.redis.pool.pool-size", 100)
	viper.SetDefault("cache.redis.pool.max-idle-conns", 10)
	viper.SetDefault("cache.redis.pool.conn-max-idle-time", 600)
	viper.SetDefault("cache.redis.pool.conn-max-life-time", 3600)
	viper.SetDefault("cache.redis.max-retries", 1)
	viper.SetDefault("cache.redis.pool-timeout", 10)
	viper.SetDefault("cache.redis.dial-timeout", 5)
	viper.SetDefault("cache.redis.read-timeout", 5)
	viper.SetDefault("cache.redis.write-timeout", 3)
	viper.SetDefault("remoteconfig.enable", true)
	viper.SetDefault("remoteconfig.provider", "etcd3")
	viper.SetDefault("remoteconfig.endpoint", "127.0.0.1:2379")
	viper.SetDefault("remoteconfig.path", "/config/prod/cloud/region/github.com/youngshawn/go-proect-demo/course/course.yaml")
	viper.SetDefault("remoteconfig.format", "yaml")
	viper.SetDefault("vault.address", "http://127.0.0.1:8200")
	viper.SetDefault("vault.auth.roleid-file-path", "~/.course-roleid")
	viper.SetDefault("vault.auth.secretid-file-path", "~/.course-secretid")
	viper.SetDefault("vault.auth.wrapped", false)
	viper.SetDefault("vault.transit.key", "course")
}

func ExposeConfigAsPFlags(cmd *cobra.Command) {
	pflags := cmd.PersistentFlags()
	pflags.StringP("listen", "l", "", "server address (default is ':8080')")
	pflags.String("database.type", "", "database type, sqlite or mysql (default is 'sqlite')")
	pflags.String("database.sqlite.filename", "", "sqlite db filename (default is './course.db')")
	pflags.String("database.sqlite.options", "", "sqlite options (default is '_foreign_keys=on')")
	pflags.String("database.mysql.username", "", "mysql username")
	pflags.String("database.mysql.password", "", "mysql password")
	pflags.String("database.mysql.address", "", "mysql address (default is '127.0.0.1:3306')")
	pflags.String("database.mysql.dbname", "", "mysql db name (default is 'course')")
	pflags.String("database.mysql.options", "", "mysql options (default is 'charset=utf8mb4&parseTime=True&loc=Local')")
	pflags.Uint("database.pool.max-idle-conns", 0, "database connection pool max-idle-conns (default is 10)")
	pflags.Uint("database.pool.max-open-conns", 0, "database connection pool max-open-conns (default is 100)")
	pflags.Uint("database.pool.conn-max-idle-time", 0, "database connection pool conn-max-idle-time in seconds (default is 300s)")
	pflags.Uint("database.pool.conn-max-life-time", 0, "database connection pool conn-max-life-time in seconds (default is 3600s)")
	pflags.Bool("cache.enable-redis", false, "enable redis (default is true)")
	pflags.Bool("cache.enable-local-cache", false, "enable local cache (default is true)")
	pflags.Bool("cache.enable-null-result-cache", false, "enable null-result cache (default is true)")
	pflags.Uint("cache.cache-ttl", 0, "cache ttl in seconds (default is 3600)")
	pflags.String("cache.redis.address", "", "redis address (default is '127.0.0.1:6379')")
	pflags.String("cache.redis.password", "", "redis password (default is '')")
	pflags.Uint("cache.redis.db", 0, "redis db index (default is 0)")
	pflags.Uint("cache.redis.pool.pool-size", 0, "redis connection pool size (default is 100)")
	pflags.Uint("cache.redis.pool.max-idle-conns", 0, "redis connection pool max-idle-conns (default is 10)")
	pflags.Uint("cache.redis.pool.conn-max-idle-time", 0, "redis connection pool conn-max-idle-time (default is 600)")
	pflags.Uint("cache.redis.pool.conn-max-life-time", 0, "redis connection pool conn-max-life-time (default is 3600)")
	pflags.Uint("cache.redis.max-retries", 0, "redis max-retries (default is 1)")
	pflags.Uint("cache.redis.pool-timeout", 0, "redis pool-timeout (default is 10)")
	pflags.Uint("cache.redis.dial-timeout", 0, "redis dial-timeout (default is 5)")
	pflags.Uint("cache.redis.read-timeout", 0, "redis read-timeout (default is 5)")
	pflags.Uint("cache.redis.write-timeout", 0, "redis write-timeout (default is 3)")
	pflags.Bool("remoteconfig.enable", false, "enable remote config (default is true)")
	pflags.String("remoteconfig.provider", "", "remote config provider (default is etcd3)")
	pflags.String("remoteconfig.endpoint", "", "remote config endpoint (default is '127.0.0.1:2379')")
	pflags.String("remoteconfig.path", "", "remote config path (default is '/config/prod/cloud/region/github.com/youngshawn/go-proect-demo/course/course.yaml')")
	pflags.String("remoteconfig.format", "", "remote config format (default is 'yaml')")
	pflags.String("vault.address", "", "vault address (default is 'http://127.0.0.1:8200')")
	pflags.String("vault.auth.roleid-file-path", "", "vault approle roleid file path (default is '~/.course-roleid')")
	pflags.String("vault.auth.secretid-file-path", "", "vault approle secretid file path (default is '~/.course-secretid')")
	pflags.Bool("vault.auth.wrapped", false, "if secretid is wrapped (default is false)")
	pflags.String("vault.transit.key", "", "vault transit key name (default is 'course')")
}
