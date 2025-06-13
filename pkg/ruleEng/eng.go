package ruleEng

import (
	"github.com/Knetic/govaluate"
)

var (
	FuncMap = make(map[string]govaluate.ExpressionFunction)
)

func InitFuncMap(funcName string, f govaluate.ExpressionFunction) {
	FuncMap[funcName] = f
}

func Run(rule, target string) (bool, error) {
	rule = processRule(rule)
	expr, err := govaluate.NewEvaluableExpressionWithFunctions(rule, FuncMap)
	if err != nil {
		return false, err
	}
	param := map[string]interface{}{
		"input": target,
	}
	result, err := expr.Evaluate(param)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// // processRule 将原本rule 处理成规则引擎能接受的样子
//
//	func processRule(rawRule string) string {
//		var funcNames []string
//		for k := range FuncMap {
//			funcNames = append(funcNames, k)
//		}
//		funcNamesString := strings.Join(funcNames, "|")
//		processRuleStr := fmt.Sprintf(ProcessRuleReg, funcNamesString)
//		re := regexp.MustCompile(processRuleStr)
//		return re.ReplaceAllStringFunc(rawRule, func(s string) string {
//			subMatch := re.FindStringSubmatch(s)
//			if len(subMatch) < 3 {
//				return s
//			}
//			funcName := subMatch[1]
//			params := subMatch[2]
//			if strings.TrimSpace(params) == "" {
//				return fmt.Sprintf("%s(input)", funcName)
//			} else {
//				return fmt.Sprintf("%s(%s,input)", funcName, params)
//			}
//		})
//	}
//
// processRule 将原本rule 处理成规则引擎能接受的样子
func processRule(expr string) string {
	ret := ""
	flag1 := false //用于控制"
	flag2 := false // 用于控制(
	isIn := false
	for _, v := range expr {
		strValue := string(v)
		addValue := strValue
		if strValue == "\"" {
			flag1 = !flag1
			if flag1 {
				isIn = true
			}
		} else if strValue == "(" {
			flag2 = true
		} else if strValue == ")" && flag2 && !flag1 {
			if isIn {
				addValue = ", input)"
			} else {
				addValue = "input)"
			}
			flag2 = false
			isIn = false
		}
		ret += addValue
	}
	return ret
}
