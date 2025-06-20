package config

import (
	"ScanIDOR/internal/pkg/ai"
	"ScanIDOR/pkg/logger"
)

type Config struct {
	LogConf  logger.LogConf `mapstructure:"LogConf"`
	AiConfig *ai.Config     `mapstructure:"AiConf"`
}
