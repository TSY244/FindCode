package fcFlag

import (
	"ScanIDOR/internal/config"
)

var (
	ConfigPath      string
	RulePath        string
	LogicDir        string
	OutputFile      string
	GoTarget        string
	AiCycle         int
	AiConfigPath    string
	AiMode          bool
	IsAutoFrameScan bool
	RunMode         string
)

var (
	CoreConf *config.Config
)

var (
	Env = map[string]string{}
)

//func init() {
//	//EnvMap = make(map[string]string)
//}
