package global

import "ScanIDOR/internal/pkg/msg"

var (
	FinishedTask = msg.NewMsg() // 用于后续存放已经完成的扫描任务队列，而不用一直等待
)
