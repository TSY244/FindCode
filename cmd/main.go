package main

import (
	"ScanIDOR/internal/pkg/config"
	"ScanIDOR/internal/pkg/env"
	"ScanIDOR/internal/pkg/rule"
	"ScanIDOR/internal/pkg/scanner"
	"ScanIDOR/pkg/logger"
	"ScanIDOR/pkg/utils"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		logger.Infof("项目花费时间是: %.3f", time.Since(start).Seconds())
	}()

	// 解析命令行参数，并且注册到"环境变量中"
	env.CheckFlag()

	conf, err := config.Init(env.ConfigPath)
	if err != nil {
		fmt.Println(err)
	}

	// 注册全局logger
	logger.SetDefaultLogger(logger.NewLogger(&conf.LogConf))

	// 加载rule
	var r rule.Rule
	if err := utils.LoadYaml(env.RulePath, &r); err != nil {
		logger.Fatal(err)
	}

	if err := scanner.Scan(env.LogicDir, &r); err != nil {
		logger.Fatal(err)
	}

}
