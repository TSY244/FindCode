package api

import (
	"ScanIDOR/internal/pkg/rule"
	"ScanIDOR/internal/pkg/scanner"
	"ScanIDOR/internal/pkg/util/consts"
	"ScanIDOR/pkg/logger"
	"ScanIDOR/pkg/utils/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func Scan(c *gin.Context) {
	// 接受一个git url
	// 返回一个结果html 用于展示结果
	gitUrl := c.PostForm("gitUrl")
	scanType := c.PostForm("type")

	rulePath, ok := consts.TypeMap[scanType]
	if !ok {
		var msgs []string
		for k, _ := range consts.TypeMap {
			msgs = append(msgs, k)
		}
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg":      "type 出错。type 只支持：" + strings.Join(msgs, ","),
			"go.error": "",
		})
		return
	}

	if ret := util.CheckGitUrl(gitUrl); !ret {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg":      "gitUrl 出错：" + gitUrl,
			"go.error": "",
		})
		return
	}
	clonePath, err := util.CloneRepository(c, gitUrl)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg":      "clone 失败 gitUrl:" + gitUrl,
			"go.error": err.Error(),
		})
		return
	}

	var r rule.Rule
	if err := util.LoadYaml(rulePath, &r); err != nil {
		logger.Fatal(err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg":      "rule.yaml 出错：" + rulePath,
			"go.error": "",
		})
		return
	}

	if err := scanner.Scan(clonePath, &r); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg":      "扫描失败",
			"go.error": err.Error(),
		})
		return
	}

	defer os.RemoveAll(clonePath)

	c.HTML(http.StatusOK, "results.html", gin.H{
		"msg":    "扫描成功",
		"result": scanner.Result, // 传递结果给模板
	})

}
