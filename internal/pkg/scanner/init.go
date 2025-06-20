package scanner

import (
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/pkg/ruleEng"
)

// 这里设置 注册是无感的
func init() {
	// 注册规则引擎的处理函数
	ruleEng.InitFuncMap(consts.ContainFunc, Contain)
	ruleEng.InitFuncMap(consts.BeginStrFunc, BeginStr)
	ruleEng.InitFuncMap(consts.EndStrFunc, EndStr)
	ruleEng.InitFuncMap(consts.RegFunc, Reg)
	ruleEng.InitFuncMap(consts.EqualFunc, Equal)
	ruleEng.InitFuncMap(consts.BeginWithLower, BeginWithLowerCase)
	ruleEng.InitFuncMap(consts.BeginWithUpper, BeginWithUpperCase)

	// 注册默认的task 处理逻辑
	AddFilterFunc(filterFuncName)
	AddFilterFunc(filterParamType)
	AddFilterFunc(filterParamName)
	AddFilterFunc(filterReturn)
	AddFilterFunc(filterRecvName)
	AddFilterFunc(filterRecvType)
}
