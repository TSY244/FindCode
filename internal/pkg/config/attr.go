package config

import (
	"ScanIDOR/pkg/logger"
)

type Config struct {
	LogConf logger.LogConf `mapstructure:"LogConf"`
}
