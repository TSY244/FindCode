package env

import (
	"ScanIDOR/internal/config"
	"ScanIDOR/internal/pkg/rule"
)

var (
	ConfigPath   string
	RulePath     string
	LogicDir     string
	OutputFile   string
	GoTarget     string
	AiCycle      int
	AiConfigPath string
	AiMode       bool
)

var (
	CoreRule *rule.Rule
	CoreConf *config.Config
)

var (
	Env = map[string]string{}
)

//func init() {
//	//EnvMap = make(map[string]string)
//}
