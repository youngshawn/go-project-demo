package config

type config struct {
	Database struct {
		Type   string
		Sqlite struct {
			DBname string
		}
		MySQL struct {
			Username string
			Password string
			Address  string
			DBname   string
			Options  string
		}
	}
	Cache struct {
		Enable bool
		Redis  struct {
			Address string
			DBindex uint `mapstructure:"db"`
		}
		EnableNullResultCache bool `mapstructure:"enable-null-result-cache"`
		EnableLocalCache      bool `mapstructure:"enable-local-cache"`
		CacheTTLInSeconds     uint `mapstructure:"cache-ttl-in-seconds"`
	}
}

var Config config

/*
func init() {
	// viper set defaults
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.sqlite.dbname", "course.db")
	viper.SetDefault("database.mysql.address", "127.0.0.1:3306")
	viper.SetDefault("database.mysql.dbname", "course")
	viper.SetDefault("database.mysql.options", "charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("cache.enable", "1")
	viper.SetDefault("cache.enable-null-result-cache", "1")
	viper.SetDefault("enable-local-cache", "0")
	viper.SetDefault("cache-ttl-in-seconds", "3600")
	viper.SetDefault("cache.redis.address", "127.0.0.1:6379")
	viper.SetDefault("cache.redis.db", "0")
}
*/
