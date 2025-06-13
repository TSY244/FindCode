package env

import (
	"ScanIDOR/internal/pkg/config"
	"ScanIDOR/internal/pkg/rule"
)

var (
	ConfigPath string
	RulePath   string
	LogicDir   string
)

var (
	CoreRule *rule.Rule
	CoreConf *config.Config
)
