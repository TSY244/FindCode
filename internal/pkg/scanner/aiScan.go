package scanner

import (
	"ScanIDOR/internal/pkg/ai"
	"ScanIDOR/internal/pkg/ai/prompt"
	"ScanIDOR/internal/pkg/ai/request"
	"ScanIDOR/internal/pkg/ai/respose"
	result2 "ScanIDOR/internal/pkg/ai/result"
	"ScanIDOR/internal/pkg/env"
	"ScanIDOR/internal/util/utils"
	"ScanIDOR/pkg/logger"
	"ScanIDOR/pkg/sysEnv"
	"ScanIDOR/pkg/template"
	"ScanIDOR/utils/util"
	"fmt"
	"github.com/goccy/go-json"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type AiResultUnit map[string][]Unit

type Unit struct {
	Result string
	Reason string
}

func aiScan() error {
	for path, apis := range apiCache {
		for _, api := range apis {
			repeatNum := 0
			const maxRepeatNum = 3
			funcCode, err := util.Decompress(api.Code)
			if err != nil {
				return err
			}

			allSubCode, err := getAllSubFuncCode(&api, path)
			if err != nil {
				return err
			}

			totalPrompt := getTotalPrompt(string(funcCode), allSubCode)

			r, err := getChatRequest()
			if err != nil {
				logger.Error(err)
				return err
			}

			jsonBody, err := getDeepseekRequest(r, totalPrompt)
			if err != nil {
				logger.Error(err)
				return err
			}

			//body := fmt.Sprintf(r.Body, prompt.JsonSystem, totalPrompt)
			r.Body = jsonBody
			for i := 0; i < env.AiCycle; i++ {
				var ret respose.DeepseekResp
				err = r.Send(&ret)
				if err != nil {
					logger.Error(err)
					return err
				}
				jsonData := utils.ExtractJSON(ret.GetChatContent())
				var jsonRet result2.JsonResult
				if err = yaml.Unmarshal([]byte(jsonData), &jsonRet); err != nil {
					logger.Error(err)
					continue
				}
				if jsonRet.Result != "true" && jsonRet.Result != "false" && repeatNum < maxRepeatNum {
					repeatNum++
					i -= 1
					continue
				}
				logger.Debug(jsonRet.Result)

				// 添加到ai result
				saveToAiResult(path, jsonRet, api)
			}
		}
	}
	SaveAiResult()
	return nil
}

func getAllSubFuncCode(api *cacheUnit, path string) (string, error) {
	allSubCodes, err := getAllSubCode(api.FuncAst, path)
	if err != nil {
		return "", err
	}
	var allSubCode string
	for i, code := range allSubCodes {
		if len(allSubCode) >= 500 {
			break
		}
		allSubCode += fmt.Sprintf("第%d段子调用代码如下：", i)
		allSubCode += fmt.Sprintf("%s\n\n", code)
	}
	return allSubCode, nil
}

func getTotalPrompt(funcCode, allSubCode string) string {
	totalPrompt := fmt.Sprintf(prompt.CheckApiPrompt, funcCode, allSubCode)
	totalPrompt = strings.Replace(totalPrompt, "\n", ";", -1)
	return totalPrompt
}

func getChatRequest() (*request.ChatRequest, error) {
	content, err := os.ReadFile(env.AiConfigPath)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	aiSk := sysEnv.GetEnv(ai.AiSkEnv)
	params := map[string]string{
		"ai_sk":  aiSk,
		"msg":    "",
		"system": "",
	}
	source := template.NewTemplate(string(content), params)
	source.Load()
	result, err := source.Replace()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var r request.ChatRequest
	if err = yaml.Unmarshal([]byte(result), &r); err != nil {
		logger.Error(err)
		return nil, err
	}
	return &r, nil
}

func getDeepseekRequest(r *request.ChatRequest, totalPrompt string) (string, error) {
	var deepseekreq request.DeepseekReq
	if err := json.Unmarshal([]byte(r.Body), &deepseekreq); err != nil {
		logger.Error(err)
	}
	var msgs []request.DeepseekMessage
	msgs = append(msgs, request.DeepseekMessage{
		Role:    "system",
		Content: prompt.JsonSystem,
	})
	msgs = append(msgs, request.DeepseekMessage{
		Role:    "user",
		Content: totalPrompt,
	})
	deepseekreq.Messages = msgs

	jsonBody, err := json.Marshal(deepseekreq)
	if err != nil {
		logger.Error(err)
	}
	return string(jsonBody), nil
}

func saveToAiResult(path string, jsonRet result2.JsonResult, api cacheUnit) {
	if resultMap, ok := AiResult[path]; ok {
		if units, ok2 := resultMap[api.FuncAst.Name.Name]; ok2 {
			units = append(units, Unit{
				Result: jsonRet.Result,
				Reason: jsonRet.Reason,
			})
			resultMap[api.FuncAst.Name.Name] = units
		}
	} else {
		AiResult[path] = map[string][]Unit{
			api.FuncAst.Name.Name: {
				Unit{
					Result: jsonRet.Result,
					Reason: jsonRet.Reason,
				},
			},
		}
	}

}
