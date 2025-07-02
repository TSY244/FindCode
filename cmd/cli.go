package main

import (
	"ScanIDOR/internal/config"
	"ScanIDOR/internal/pkg/fcFlag"
	"ScanIDOR/internal/pkg/rule"
	"ScanIDOR/internal/pkg/scanner"
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/pkg/fingerprint"
	"ScanIDOR/pkg/logger"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		logger.Infof("项目花费时间是: %.3f", time.Since(start).Seconds())
	}()

	// 解析命令行参数，并且注册到"环境变量中"
	fcFlag.CheckFlag()

	conf, err := config.Init(fcFlag.ConfigPath)
	if err != nil {
		fmt.Println(err)
	}

	// 注册全局logger
	logger.SetDefaultLogger(logger.NewLogger(&conf.LogConf))

	// 加载rule
	if fcFlag.IsAutoFrameScan {
		for _, r := range rule.LoadRuleWithFrame(fingerprint.GetProductPrint(fcFlag.LogicDir)) {
			loadEnv(&r, conf) // 加载环境变量
			envData := scanner.NewEnv()
			envData.AiCycle = fcFlag.AiCycle
			if err := scanner.Scan(fcFlag.LogicDir, &r, envData); err != nil {
				logger.Fatal(err)
			}
		}
	} else {

	}

}

func loadEnv(r *rule.Rule, c *config.Config) {
	if fcFlag.GoTarget != "" {
		r.GoModeTargetRule.Rule = fcFlag.GoTarget
	}
	if fcFlag.AiMode == true {
		r.Mode = append(r.Mode, consts.AiMode)
		r.AiConfig = c.AiConfig
	}
}
