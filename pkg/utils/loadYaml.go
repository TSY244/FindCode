package utils

import (
	"github.com/spf13/viper"
	"path/filepath"
)

// LoadYaml 加载yaml 返回T
func LoadYaml[T any](filePath string, data *T) error {
	// 处理获取filePath
	//filePath = strings.Replace(filePath, "\\", "/", -1)
	path := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)

	// viper 配置
	viper.SetConfigName(fileName)
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")

	// 解析yaml
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(data); err != nil {
		return err
	}
	return nil
}
