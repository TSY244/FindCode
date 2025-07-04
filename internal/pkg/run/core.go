package run

import (
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/pkg/logger"
	"strings"
)

var (
	runMap = map[string]func(){
		consts.Command: Command,
		consts.Server:  Server,
	}
)

func Run(mode string) {
	mode = strings.ToLower(mode)
	logger.Info("您已启动 " + mode + " 模式")
	runMap[mode]()
}
