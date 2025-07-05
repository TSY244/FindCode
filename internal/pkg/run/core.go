package run

import (
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/pkg/logger"
	"strings"
)

var (
	runMap = map[string]func(){
		consts.Command: Command, // 启动命令行模式
		consts.Server:  Server,  // 启动server 模式，使用gin 搭建的后端和前端
	}
)

func Run(mode string) {
	mode = strings.ToLower(mode)
	logger.Info("您已启动 " + mode + " 模式")
	runMap[mode]()
}
