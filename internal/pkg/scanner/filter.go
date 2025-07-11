package scanner

import (
	"ScanIDOR/internal/pkg/rule"
	"ScanIDOR/pkg/logger"
	"errors"
	"fmt"
	"go/ast"
)

type FilterFunc func(decl *ast.FuncDecl, rule *rule.FuncRuleUnit) (bool, error)

type FilterFuncs []FilterFunc

var (
	filterFuncs = make(FilterFuncs, 0)
)

func AddFilterFunc(taskFunc FilterFunc) {
	filterFuncs = append(filterFuncs, taskFunc)
}

func filter(decl *ast.FuncDecl, rule *rule.Rule) (bool, error) {
	for _, funcRule := range rule.FuncRules { // 遍历每一个func rule
		for i, filterFunc := range filterFuncs {
			ok, err := filterFunc(decl, &funcRule)
			if err != nil {
				return false, err
			}
			if !ok {
				break
			}
			if i == len(filterFuncs)-1 {
				return true, nil
			}
		}
	}
	return false, nil
}

func filterFuncName(decl *ast.FuncDecl, funcRule *rule.FuncRuleUnit) (bool, error) {
	if funcRule.FuncNameRule == nil {
		return true, nil
	}
	if funcRule.FuncNameRule.Rule == "" {
		return false, nil // 控制循环，并非真false
	}

	if ret, err := matchStr(funcRule.FuncNameRule.Rule, decl.Name.Name); err != nil {
		logger.Error(err)
	} else if ret {
		return true, nil
	}
	return false, nil

}

func filterParamType(decl *ast.FuncDecl, funcRule *rule.FuncRuleUnit) (bool, error) {

	if funcRule.ParamTypeRule == nil || funcRule.ParamTypeRule.Rule == nil {
		return true, nil
	}
	if len(funcRule.ParamTypeRule.Rule) == 0 {
		return true, nil
	}
	params, err := getFuncParamTypes(decl)
	if err != nil {
		return false, err
	}
	for _, paramUnit := range funcRule.ParamTypeRule.Rule {
		if paramUnit.Size == -1 {
			return true, nil
		}
		if len(params) != paramUnit.Size {
			continue
		}
		if paramUnit.Rules != nil {
			if ret, err := matchStrSplice(params, paramUnit.Rules); err != nil {
				return false, err
			} else if ret {
				return true, nil
			}
		} else { // 没有规则直接返回
			return true, nil
		}
	}
	return false, nil

}

func filterParamName(decl *ast.FuncDecl, funcRule *rule.FuncRuleUnit) (bool, error) {
	if funcRule.ParamNameRule == nil || funcRule.ParamNameRule.Rule == nil {
		return true, nil
	}
	if len(funcRule.ParamNameRule.Rule) == 0 {
		return true, nil
	}
	names, err := getFuncParamNames(decl)
	if err != nil {
		return false, err
	}
	for _, paramUnit := range funcRule.ParamNameRule.Rule {
		if paramUnit.Size == -1 {
			return true, nil
		}
		if len(names) != paramUnit.Size {
			continue
		}
		if paramUnit.Rules != nil {
			if ret, err := matchStrSplice(names, paramUnit.Rules); err != nil {
				return false, err
			} else if ret {
				return true, nil
			}
		} else {
			return true, nil
		}
	}
	return false, nil

}

func getFuncParamTypes(funcDecl *ast.FuncDecl) ([]string, error) {
	var funcParamTypes []string
	for _, param := range funcDecl.Type.Params.List {
		funcType, err := getParamType(param)
		if err != nil {
			return nil, err
		}
		funcParamTypes = append(funcParamTypes, funcType)
	}
	return funcParamTypes, nil
}

func getParamType(param *ast.Field) (string, error) {
	var funcParamTypes string
	if pType, ok := param.Type.(*ast.SelectorExpr); ok {
		if X, ok := pType.X.(*ast.Ident); ok {
			return X.Name + "." + pType.Sel.Name, nil
		}
	}
	// 针对项目特性，没有深度解析
	if pType, ok := param.Type.(*ast.StarExpr); ok {
		if X, ok := pType.X.(*ast.SelectorExpr); ok {
			if X2, ok := X.X.(*ast.Ident); ok {
				return "*" + X2.Name + "." + X.Sel.Name, nil
			}
		}
	}
	return funcParamTypes, nil
}

func getFuncParamNames(funcDecl *ast.FuncDecl) ([]string, error) {
	var funcParamNames []string
	for _, param := range funcDecl.Type.Params.List {
		funcType, err := getParamNames(param)
		if err != nil {
			return nil, err
		}
		funcParamNames = append(funcParamNames, funcType...)
	}
	return funcParamNames, nil
}

func getParamNames(param *ast.Field) ([]string, error) {
	names := make([]string, 0)
	for _, astName := range param.Names {
		names = append(names, astName.Name)
	}
	return names, nil
}

// 1. size
// 2. return type

func filterReturn(decl *ast.FuncDecl, funcRule *rule.FuncRuleUnit) (bool, error) {
	if funcRule.ReturnTypeRule == nil || funcRule.ReturnTypeRule.Rules == nil {
		return true, nil
	}
	if len(funcRule.ReturnTypeRule.Rules) == 0 {
		return true, nil
	}
	for _, unit := range funcRule.ReturnTypeRule.Rules {
		types, err := getReturnTypes(decl)
		if err != nil {
			return false, err
		}
		if unit.Size == -1 {
			return true, nil
		}

		if len(types) != unit.Size {
			continue
		}
		if unit.Rules != nil {
			if ret, err := matchStrSplice(types, unit.Rules); err != nil {
				return false, err
			} else if ret {
				return true, nil
			}
		} else {
			return true, nil
		}
	}
	return false, nil
}

// getReturnTypes 获取返回值
func getReturnTypes(decl *ast.FuncDecl) ([]string, error) {
	var funcReturnTypes []string
	if decl.Type.Results == nil {
		return funcReturnTypes, nil
	}
	for _, field := range decl.Type.Results.List {
		typeStr := exprToString(field.Type)
		funcReturnTypes = append(funcReturnTypes, typeStr)
	}
	return funcReturnTypes, nil
}

// exprToString 解析返回值
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.FuncType:
		return "func"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		return "chan " + exprToString(t.Value)
	case *ast.Ellipsis:
		return "..." + exprToString(t.Elt)
	case *ast.StructType:
		return "struct{...}"
	case *ast.ParenExpr:
		return "(" + exprToString(t.X) + ")"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func checkRecv(decl *ast.FuncDecl) error {
	if decl.Recv == nil {
		return NoRecv
	}
	if len(decl.Recv.List) != 1 {
		return RecvErr
	}
	return nil
}

func filterRecvName(decl *ast.FuncDecl, funcRule *rule.FuncRuleUnit) (bool, error) {
	if funcRule.RecvNameRule == nil {
		return true, nil
	}
	if funcRule.RecvNameRule.Rule == "" {
		return false, nil // 控制循环，并非真false
	}

	recvName, err := getRecvName(decl)
	if err != nil {
		if errors.Is(err, NoRecv) {
			return true, nil
		}
		return false, err
	}
	if ret, err := matchStr(funcRule.RecvNameRule.Rule, recvName); err != nil {
		logger.Error(err)
	} else if ret {
		return true, nil
	}
	return false, nil
}

func filterRecvType(decl *ast.FuncDecl, funcRule *rule.FuncRuleUnit) (bool, error) {
	if funcRule.RecvTypeRule == nil {
		return true, nil
	}
	if funcRule.RecvTypeRule.Rule == "" {
		return false, nil // 控制循环，并非真false
	}

	recvType, err := getRecvType(decl)
	if err != nil {
		if errors.Is(err, NoRecv) {
			return true, nil
		}
		return false, err
	}
	if ret, err := matchStr(funcRule.RecvTypeRule.Rule, recvType); err != nil {
		logger.Error(err)
	} else if ret {
		return true, nil
	}
	return false, nil
}

func getRecvName(decl *ast.FuncDecl) (string, error) {
	if err := checkRecv(decl); err != nil {
		return "", err
	}
	if len(decl.Recv.List[0].Names) == 0 {
		return "", errors.New("receiver error")
	}
	return decl.Recv.List[0].Names[0].Name, nil // go 默认值会有一个接受值
}

func getRecvType(decl *ast.FuncDecl) (string, error) {
	if err := checkRecv(decl); err != nil {
		return "", err
	}
	recvType := decl.Recv.List[0].Type
	if pType, ok := recvType.(*ast.SelectorExpr); ok {
		if X, ok := pType.X.(*ast.Ident); ok {
			return X.Name + "." + pType.Sel.Name, nil
		}
	}
	// 针对项目特性，没有深度解析
	if pType, ok := recvType.(*ast.StarExpr); ok {
		if X, ok := pType.X.(*ast.SelectorExpr); ok {
			if X2, ok := X.X.(*ast.Ident); ok {
				return "*" + X2.Name + "." + X.Sel.Name, nil
			}

		}
		if X, ok := pType.X.(*ast.Ident); ok {
			return "*" + X.Name, nil
		}
	}
	return "", errors.New("receiver error")
}
