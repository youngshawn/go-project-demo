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
		Redis struct {
			Address string
			DB      uint
		}
		EnableNullResultCache bool
		EnableLocalCache      bool
		CacheTTLInSeconds     uint
	}
}
