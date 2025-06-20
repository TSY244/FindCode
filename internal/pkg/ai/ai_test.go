package ai

import (
	"ScanIDOR/internal/config"
	"ScanIDOR/internal/pkg/ai/prompt"
	"ScanIDOR/internal/pkg/ai/request"
	"ScanIDOR/internal/pkg/ai/respose"
	result2 "ScanIDOR/internal/pkg/ai/result"
	"ScanIDOR/internal/util/utils"
	"ScanIDOR/pkg/logger"
	"ScanIDOR/pkg/sysEnv"
	"ScanIDOR/pkg/template"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestRequest(t *testing.T) {

	configPath := "/Users/august/code/go/FindCode/etc/config.yaml"
	c, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	aiSk := sysEnv.GetEnv("ai_sk")
	params := map[string]string{
		"env.api_sk": aiSk,
		"msg":        "请你计算1+1 是否等于2，结果请用纯json 字符串。不带md 格式",
		"system":     prompt.CheckApiSystem,
	}
	source := template.NewTemplate(string(content), params)
	source.Load()
	result, err := source.Replace()
	if err != nil {
		logger.Error(err)
		return
	}

	var r request.ChatRequest
	if err = yaml.Unmarshal([]byte(result), &r); err != nil {
		logger.Error(err)
		return
	}

	var ret respose.DeepseekResp
	err = r.Send(&ret)
	if err != nil {
		logger.Fatal(err)
	}
	jsonData := utils.ExtractJSON(ret.GetChatContent())
	var jsonRet result2.JsonResult
	if err = yaml.Unmarshal([]byte(jsonData), &jsonRet); err != nil {
		logger.Error(err)
		return
	}
	fmt.Println(jsonRet)
}
