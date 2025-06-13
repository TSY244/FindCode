package scanner

import (
	"ScanIDOR/pkg/logger"
	"bufio"
	"fmt"
	"os"
	"time"
)

func SaveToFile(taskName string) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	fileName := taskName + "_" + timestamp + ".txt"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("create file %s failed.\n", fileName)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	var msg string
	if strModeFlag {
		msg = "存在filter\n"
	} else {
		msg = "不存在filter\n"
	}
	_, err = writer.WriteString(msg)
	if err != nil {
		logger.Error(err)
	}

	for path, ret := range Result {
		strMeg := fmt.Sprintf("\n\n%s 中以下的函数可能存在越权漏洞\n开始行数:结尾行数:函数名字\n", path)
		_, err := writer.WriteString(strMeg)
		if err != nil {
			logger.Errorf("save file %s failed.\n", fileName)
			return
		}
		for _, v := range ret {
			_, err := writer.WriteString(v + "\n")
			if err != nil {
				logger.Errorf("save file %s failed.\n", fileName)
				return
			}
		}
	}
	logger.Infof("save file %s success.", fileName)
}
