package controllers

import (
	"ScanIDOR/internal/pkg/rule"
	"ScanIDOR/internal/pkg/scanner"
	"ScanIDOR/internal/pkg/server/dtos/requests"
	"ScanIDOR/internal/pkg/server/services"
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/internal/util/utils"
	"ScanIDOR/pkg/fingerprint"
	"ScanIDOR/pkg/logger"
	util2 "ScanIDOR/utils/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

type FindCodeController struct {
	findCodeService services.FindCodeService
}

func NewFindCodeController(findCodeService services.FindCodeService) *FindCodeController {
	return &FindCodeController{findCodeService: findCodeService}
}

func (f *FindCodeController) Scan(c *gin.Context) {
	// 接受一个git url
	// 返回一个结果html 用于展示结果
	var request requests.APIScanRequest
	if err := c.ShouldBind(&request); err != nil {
		// 若绑定失败，返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "表单数据解析失败",
			"details": err.Error(),
		})
		return
	}

	gitUrl := request.GitURL
	//scanType := request.Type
	isUseAi := request.IsUseAi
	authenticationCodes := request.AuthenticationCodes
	temp := make([]string, 0)
	for _, str := range authenticationCodes {
		if "" == str || "[]" == str {
			continue
		}
		temp = append(temp, str)
	}
	authenticationCodes = temp

	if ret := util2.CheckGitUrl(gitUrl); !ret {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg.go":   "gitUrl 出错：" + gitUrl,
			"go.error": "",
		})
		return
	}
	clonePath, err := util2.CloneRepository(c, gitUrl)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg.go":   "clone 失败 gitUrl:" + gitUrl,
			"go.error": err.Error(),
		})
		return
	}
	time.Sleep(2 * time.Second) // 给本地加载文件留时间

	scanType := fingerprint.GetProductPrint(clonePath)
	rulePath, ok := consts.TypeMap[scanType[0]]
	if !ok {
		var msgs []string
		for k, _ := range consts.TypeMap {
			msgs = append(msgs, k)
		}
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg.go":   "type 出错。type 只支持：" + strings.Join(msgs, ","),
			"go.error": "",
		})
		return
	}

	defer func() {
		os.RemoveAll(clonePath)
	}()

	var r rule.Rule
	if err := util2.LoadYaml(rulePath, &r); err != nil {
		logger.Fatal(err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg.go":   "rule.yaml 出错：" + rulePath,
			"go.error": "",
		})
		return
	}
	//ctx := context.WithValue(context.Background(), consts.IsUseCtxKey, false)

	if isUseAi {
		for i, m := range r.Mode {
			if m == consts.AiMode {
				break
			}
			if len(r.Mode)-1 == i {
				r.Mode = append(r.Mode, consts.AiMode)
			}
		}
		aiConfig, err := utils.GetAiConfig(request.Model)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"msg.go":   "aiConfig 出错：" + rulePath,
				"go.error": "",
			})
			return
		}
		r.AiConfig = aiConfig
		r.AiConfig.IsUseAiPrompt = request.IsUseAiPrompt
		if request.IsUseAiPrompt {
			r.AiConfig.IsReturnBool = request.IsReturnBool
		}
		r.AiConfig.Prompt = request.AiPrompt
	}

	env := scanner.NewEnv()
	env.AiCycle = consts.ServerMaxCycle

	if len(authenticationCodes) != 0 {
		r.GoModeTargetRule.Rule = scanner.GetContainsRule(authenticationCodes)
	}

	if err := scanner.Scan(clonePath, &r, env); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg.go":   "扫描失败",
			"go.error": err.Error(),
		})
		return
	}

	ret, boolRet := scanner.GetResult(clonePath, env)

	if isUseAi {
		c.HTML(http.StatusOK, "ai_results.html", gin.H{
			"msg.go":       "扫描成功",
			"result":       boolRet,
			"isReturnBool": !request.IsUseAiPrompt || request.IsReturnBool,
		})
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"msg.go": "扫描成功",
		"result": ret, // 传递结果给模板
	})

}

func (f *FindCodeController) ShowScanHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "scan_code.html", nil)
}
