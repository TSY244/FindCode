package scanner

import (
	"ScanIDOR/internal/pkg/env"
	"ScanIDOR/pkg/logger"
	"bufio"
	"fmt"
	"os"
	"time"
)

func SaveToFile(taskName string) {
	resultDir := "result/"
	err := os.MkdirAll(resultDir, os.ModePerm)
	if err != nil {
		return
	}
	var fileName string

	if env.OutputFile != "" {
		fileName = resultDir + env.OutputFile
	} else {
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		fileName = resultDir + taskName + "_" + timestamp + ".txt"
	}

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
		strMeg := fmt.Sprintf("\n\n%s 扫描结果如下\n开始行数:结尾行数:函数名字\n", path)
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

func SaveAiResult() {
	resultDir := "ai_result/"
	err := os.MkdirAll(resultDir, os.ModePerm)
	if err != nil {
		return
	}
	for path, aiResult := range AiResult {
		for funcName, units := range aiResult {
			timestamp := fmt.Sprintf("%d", time.Now().Unix())
			fileName := resultDir + funcName + "_" + timestamp + ".txt"
			file, err := os.Create(fileName)
			if err != nil {
				logger.Errorf("save file %s failed.\n", fileName)
			}
			defer file.Close()
			writer := bufio.NewWriter(file)
			defer writer.Flush()
			var msg string
			msg = path + " 结果如下:\n" + "funcName: " + funcName + "\n"
			for i, unit := range units {
				prefix := fmt.Sprintf("ai第%d次结果: \n", i)
				msg += prefix + "result: " + unit.Result + "\n" + "reason: " + unit.Reason + "\n\n"
			}
			_, err = writer.WriteString(msg)
			if err != nil {
				logger.Errorf("save file %s failed.\n", fileName)
			}
		}
	}
}
