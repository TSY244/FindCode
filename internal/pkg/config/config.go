package config

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var (
	CoreConfig *Config
)

func Init(path string) (*Config, error) {

	fileName := filepath.Base(path)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	filePath := filepath.Dir(path)

	// 使用viper
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filePath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&CoreConfig)
	if err != nil {
		return nil, err
	}

	return CoreConfig, nil
}
