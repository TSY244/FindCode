package scanner

import "ScanIDOR/pkg/ruleEng"

// 这里设置 注册是无感的
func init() {
	// 注册规则引擎的处理函数
	ruleEng.InitFuncMap(ContainFunc, Contain)
	ruleEng.InitFuncMap(BeginStrFunc, BeginStr)
	ruleEng.InitFuncMap(EndStrFunc, EndStr)
	ruleEng.InitFuncMap(RegFunc, Reg)
	ruleEng.InitFuncMap(EqualFunc, Equal)
	ruleEng.InitFuncMap(BeginWithLower, BeginWithLowerCase)
	ruleEng.InitFuncMap(BeginWithUpper, BeginWithUpperCase)

	// 注册默认的task 处理逻辑
	AddFilterFunc(filterFuncName)
	AddFilterFunc(filterParamType)
	AddFilterFunc(filterParamName)
	AddFilterFunc(filterReturn)
	AddFilterFunc(filterRecvName)
	AddFilterFunc(filterRecvType)
}
