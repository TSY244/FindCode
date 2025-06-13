package logger

type LogConf struct {
	FileName       string `mapstructure:"FileName"`
	MaxSize        int    `mapstructure:"MaxSize"`
	MaxBackups     int    `mapstructure:"MaxBackups"`
	MaxAge         int    `mapstructure:"MaxAge"`
	Compress       bool   `mapstructure:"Compress"`
	Level          string `mapstructure:"Level"`
	PrintToConsole bool   `mapstructure:"PrintToConsole"`
}
