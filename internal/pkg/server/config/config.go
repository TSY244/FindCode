package config

type Config struct {
	DbConfig  DBConfig
	GinConfig *GinConfig
}

func NewConfig() *Config {
	return &Config{
		DbConfig:  NewDBConfig(),
		GinConfig: NewGinConfig(),
	}
}
