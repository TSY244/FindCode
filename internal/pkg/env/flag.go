package env

import (
	"flag"
	"fmt"
)

func init() {
	flag.StringVar(&ConfigPath, "c", "./etc/config.yaml", "-c ./etc/config.yaml")
	flag.StringVar(&RulePath, "r", "", "-r ./etc/rule.yaml")
	flag.StringVar(&LogicDir, "l", "./logic/", "-l ./logic/")
	flag.StringVar(&OutputFile, "o", "", "-o ./output/")
	flag.StringVar(&GoTarget, "go_target", "", "-t ./")
	flag.IntVar(&AiCycle, "ai_cycle", 3, "-ai_cycle 3")
	flag.StringVar(&AiConfigPath, "ai_config", "./etc/deepseekConfig.yaml", "-ai_config ./etc/deepseekConfig.yaml")
	flag.BoolVar(&AiMode, "ai", false, "-ai true")
	flag.BoolVar(&IsAutoFrameScan, "auto_frame", true, "-auto_frame true")
	flag.Parse()
}

func CheckFlag() {
	if ConfigPath == "" {
		fmt.Println("you need to set config path")
		Help()
	}
	if RulePath == "" {
		IsAutoFrameScan = true
	} else {
		IsAutoFrameScan = false
	}
	if LogicDir == "./logic/" {
		fmt.Println("you need to set logic dir")
	}
}

func Help() {
	flag.Usage()
}
