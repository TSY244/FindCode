package main

import (
	"ScanIDOR/internal/config"
	"ScanIDOR/internal/pkg/fcFlag"
	"ScanIDOR/internal/pkg/global"
	"ScanIDOR/internal/pkg/rule"
	"ScanIDOR/internal/pkg/scanner"
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

	fcFlag.CheckFlag()

	conf, err := config.Init(fcFlag.ConfigPath)
	if err != nil {
		fmt.Println(err)
	}

	// 注册全局logger
	logger.SetDefaultLogger(logger.NewLogger(&conf.LogConf))

	// 加载rule
	rulePath := fcFlag.RulePath
	if fcFlag.IsAutoFrameScan { // 支持自动框架扫描
		// 判断使用的框架是什么
		finger, err := fingerprint.GetProductPrint(fcFlag.LogicDir)
		if err != nil {
			logger.Fatal(err.Error())
		}
		rulePath = global.RuleMap[finger]
	}

	r := rule.LoadRule(rulePath, conf)

	if err := scanner.Scan(fcFlag.LogicDir, &r, scanner.NewEnv()); err != nil {
		logger.Fatal(err)
	}
}
