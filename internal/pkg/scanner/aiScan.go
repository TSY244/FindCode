package scanner

import (
	"ScanIDOR/internal/pkg/ai"
	"ScanIDOR/internal/pkg/ai/prompt"
	"ScanIDOR/internal/pkg/ai/respose"
	result2 "ScanIDOR/internal/pkg/ai/result"
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/internal/util/utils"
	"ScanIDOR/pkg/logger"
	"ScanIDOR/utils/util"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

type AiBoolResultUnit map[string][]aiBoolUnit

type AiBoolResultUnitWithStatue map[string]aiBoolUnitWithStatue
type aiBoolUnitWithStatue struct {
	Statue      int // -1 危险 0 可疑 1 安全
	AiBoolUnits []aiBoolUnit
}

type AiResultUnit map[string][]aiStrUnit

type aiBoolUnit struct {
	Result string
	Reason string
}

type aiStrUnit struct {
	Result string
}

func aiScan(config *ai.Config, env2 *Env) error {
	for path, apis := range env2.ApiCache {
		for _, api := range apis {
			repeatNum := 0
			const maxRepeatNum = 3
			funcCode, err := util.Decompress(api.Code)
			if err != nil {
				return err
			}

			allSubCode, err := getAllSubFuncCode(&api, path, env2)
			if err != nil {
				return err
			}

			totalPrompt := getTotalPrompt(config, string(funcCode), allSubCode)

			r := utils.GetChatRequest(config)

			jsonBody, err := utils.GetDeepseekRequest(r, totalPrompt, config)
			if err != nil {
				logger.Error(err)
				return err
			}

			//body := fmt.Sprintf(r.Body, prompt.CheckApiSystem, totalPrompt)
			r.Body = jsonBody
			for i := 0; i < env2.AiCycle; i++ {
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
				saveToAiResult(path, jsonRet, api, config, env2)
			}
		}
	}
	SaveAiResult(env2)
	return nil
}

func getAllSubFuncCode(api *cacheUnit, path string, env *Env) (string, error) {
	allSubCodes, err := getAllSubCodeWithLevel(api.FuncAst, path, env, consts.FirstLevel, consts.MaxLevel)
	if err != nil {
		return "", err
	}
	var allSubCode string
	for i, code := range allSubCodes {
		if len(allSubCode) >= 40000 {
			break
		}
		allSubCode += fmt.Sprintf("第%d段子调用代码如下：", i+1)
		allSubCode += fmt.Sprintf("%s\n\n", code)
	}
	return allSubCode, nil
}

func getTotalPrompt(config *ai.Config, funcCode, allSubCode string) string {
	var totalPrompt string
	if config.IsUseAiPrompt {
		// 计算%s 字符串的个数
		size := strings.Count(config.Prompt, "%s")
		totalPrompt = config.Prompt
		switch size {
		case 0:
			//return config.Prompt
			totalPrompt = strings.Replace(totalPrompt, "\n", ";", -1)
		case 1:
			totalPrompt = fmt.Sprintf(config.Prompt, funcCode)
			totalPrompt = strings.Replace(totalPrompt, "\n", ";", -1)
		case 2:
			totalPrompt = fmt.Sprintf(config.Prompt, funcCode, allSubCode)
			totalPrompt = strings.Replace(totalPrompt, "\n", ";", -1)
		default:
			totalPrompt = fmt.Sprintf(prompt.CheckApiPrompt, funcCode, allSubCode)
		}

		if config.IsReturnBool {
			totalPrompt += prompt.ReturnBoolPrompt
		}
	} else {
		totalPrompt = fmt.Sprintf(prompt.CheckApiPrompt, funcCode, allSubCode)
		totalPrompt = strings.Replace(totalPrompt, "\n", ";", -1)
	}
	return totalPrompt
}

func saveToAiResult(path string, jsonRet result2.JsonResult, api cacheUnit, config *ai.Config, env *Env) {
	if resultMap, ok := env.AiBoolResult[path]; ok {
		if units, ok2 := resultMap[api.FuncAst.Name.Name]; ok2 {
			units = append(units, aiBoolUnit{
				Result: jsonRet.Result,
				Reason: jsonRet.Reason,
			})
			resultMap[api.FuncAst.Name.Name] = units
		} else {
			resultMap[api.FuncAst.Name.Name] = []aiBoolUnit{
				{
					Result: jsonRet.Result,
					Reason: jsonRet.Reason,
				},
			}
		}
	} else {
		env.AiBoolResult[path] = map[string][]aiBoolUnit{
			api.FuncAst.Name.Name: {
				aiBoolUnit{
					Result: jsonRet.Result,
					Reason: jsonRet.Reason,
				},
			},
		}
	}

}
