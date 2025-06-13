package scanner

import (
	"ScanIDOR/internal/pkg/rule"
	"ScanIDOR/pkg/logger"
	"ScanIDOR/pkg/utils"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

type modeFunc func(path string, target *rule.Rule) error

var (
	ModeFuncMap = map[string]modeFunc{
		GoMode:  goModeFunc,
		StrMode: strModeFunc,
	}
	strModeFlag  = false
	isUsedGoMode = false
)

func Scan(path string, target *rule.Rule) error {
	if target == nil {
		return errors.New("target is nil")
	}
	info, err := checkFileStatue(path)
	if err != nil {
		return err
	}
	// 根据文件类型启动不同的扫描逻辑
	if info.IsDir() {
		if err := dealDir(path, target, BeginLevel); err != nil {
			return err
		}
	} else {
		if err := dealFile(path, target); err != nil {
			return err
		}
	}
	// 发现存在go mode 才会进行go 的扫描
	if ret := isUseMode(target); ret {
		isUsedGoMode = true
		if err := processFuncDecls(target); err != nil {
			return err
		}
	}
	// 在下面对每一个 ast.FuncDecl 处理
	//if err := processFuncDecls(target); err != nil {
	//	return err
	//}
	SaveToFile(target.TaskName)
	printResult(Result)
	logger.Info("模块已经扫描结束")
	return nil
}

func isUseMode(target *rule.Rule) bool {
	for _, r := range target.Mode {
		if r == GoMode {
			return true
		}
	}
	return false
}

// dealDir 将目录下的文件逐个遍历，将其中的函数保存到对应的cache 中
func dealDir(path string, target *rule.Rule, deep int) error {
	files, err := getAllFilesInDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			// 通过规则可以限定对最多解析的层数
			if deep > target.Path.DeepSize {
				continue
			}
			if err := dealDir(filepath.Join(path, file.Name()), target, deep+1); err != nil {
				return err
			}
		} else {
			filePath := filepath.Join(path, file.Name())
			//if target.Path.Rule != "" {
			//	if ret, err := matchStr(target.Path.Rule, filePath); err != nil {
			//		return err
			//	} else if !ret {
			//		continue
			//	}
			//}
			if err := dealFile(filePath, target); err != nil {
				return err
			}
		}
	}
	return nil
}

// dealFile 解析文件，提取出函数，保存到cache 中
func dealFile(path string, target *rule.Rule) error {
	info, err := checkFileStatue(path)
	if err != nil {
		return err
	}
	if !isGoFile(info) {
		return nil
	}

	for _, mode := range target.Mode {
		f, ok := ModeFuncMap[mode]
		if !ok {
			logger.Warn("不存在mode: ", mode)
			continue
		}
		if err := f(path, target); err != nil {
			return err
		}
	}
	return nil
}

func strModeFunc(path string, r *rule.Rule) error {
	if strModeFlag {
		return nil
	}
	fileStr, err := readFile2String(path)
	if err != nil {
		return err
	}
	if r.StrModeTargetRule.Rule != "" {
		if ret, err := matchStr(r.StrModeTargetRule.Rule, fileStr); err != nil {
			logger.Errorf(err.Error())
			return err
		} else if !ret {
			return nil
		} else {
			strModeFlag = true
		}
	}
	return nil
}

func goModeFunc(path string, target *rule.Rule) error {

	err := processFile(path, target)
	if err != nil {
		return err
	}
	return nil
}

// processFile 解析文件真正的函数
func processFile(filePath string, rule *rule.Rule) error {
	srcStr, err := readFile2String(filePath)
	if err != nil {
		return err
	}

	// 解析文件的默认步骤
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, filePath, srcStr, 0)
	if err != nil {
		return err
	}

	//func 定义了对函数的处理逻辑
	scanDecls(fileAst.Decls, func(decl *ast.FuncDecl) {
		// 处理的是，将api func 和 非api func 补充到chache
		ret, err := filter(decl, rule)

		if !errors.Is(err, ArgsSizeNotEqualErr) && !errors.Is(err, FuncParamsNotEqualErr) && err != nil {
			return
		}
		if rule.Path.Rule != "" {
			if subRet, err := matchStr(rule.Path.Rule, filePath); err != nil {
				return
			} else if !subRet {
				ret = false
			}
		}
		if ret {
			if err := saveApiFunc(filePath, srcStr, decl, fset); err != nil {
				logger.Error(err.Error())
				return
			}
		} else {
			savaNoApiFunc(decl, srcStr, fset)
		}
	})
	return nil
}

func processFuncDecls(rule *rule.Rule) error {
	for path, apis := range apiCache {
		info, err := checkFileStatue(path)
		if err != nil {
			return err
		}
		if rule.File.Rule != "" {
			if ret, err := matchStr(rule.File.Rule, info.Name()); err != nil {
				return err
			} else if !ret {
				return nil
			}
		}

		var result []string
		for _, api := range apis {
			ret, err := processApi(api, rule, path)
			if err != nil {
				return err
			}
			result = append(result, ret...)
		}
		if len(result) > 0 {
			Result[path] = result
		}
	}
	return nil
}

func processApi(api cacheUnit, rule *rule.Rule, path string) ([]string, error) {
	//debug
	idors := make([]string, 0)
	funcCode, err := utils.Decompress(api.Code)
	if err != nil {
		return nil, err
	}

	if ret, err := matchStr(rule.GoModeTargetRule.Rule, string(funcCode)); err != nil {
		logger.Error(err.Error())
	} else if ret { // 处理子调用逻辑
		subRet, err := processFuncDecl(path, api.FuncAst, rule)
		if err != nil {
			return nil, err
		}
		if subRet {
			// 可能存在问题的，统一保存
			hashValue := GetFuncAstHash(api.FuncAst)
			JudgedCache[hashValue] = true
			startLine, endLine := getStartAndEndLine(api.Fset, api.FuncAst)
			result := fmt.Sprintf("%d:%d:%s", startLine, endLine, api.FuncAst.Name.Name)
			idors = append(idors, result)
		}
	}
	return idors, nil
}

// processFuncDecl 返回true 表示没有鉴权框架
func processFuncDecl(path string, decl *ast.FuncDecl, rule *rule.Rule) (bool, error) {
	// 1. 获取所有的子调用的函数名字
	allSubFuncDecls, names := getAllSubFuncDecls(decl, path)

	// 2. 无法获取funcdecl 的子调用，通过name 进行判断
	// todo： 和funcdecl 分开统计
	for _, name := range names {
		if ret, err := processNameDecl(name, rule, path); err != nil {
			return false, err
		} else if !ret {
			return false, nil
		}
	}

	// 统计funcdecl 的调用逻辑
	for _, subFuncDecl := range allSubFuncDecls {
		if ret, err := processSubFuncDecl(subFuncDecl, rule, path); err != nil {
			return false, err
		} else if !ret {
			return false, nil
		}
	}
	return true, nil
}

// processSubFuncDecl 处理函数调用的子函数的
func processSubFuncDecl(subFuncDecl ast.FuncDecl, rule *rule.Rule, path string) (bool, error) {
	hashKey := GetFuncAstHash(&subFuncDecl)
	if cacheRet, ok := JudgedCache[hashKey]; ok {
		if cacheRet {
			return true, nil
		}
		return cacheRet, nil
	}

	unit, ok := funcCacheMap[hashKey]
	if !ok {
		if unit == nil {
			return false, nil
		}
		name := unit.FuncAst.Name.Name
		if ret, err := processNameDecl(name, rule, path); err != nil {
			return false, err
		} else if !ret {
			return false, nil
		}
		return true, nil
	}

	funcCode := unit.Code
	if ret, err := matchStr(rule.GoModeTargetRule.Rule, string(funcCode)); err != nil {
		logger.Error(err.Error())
	} else if ret {
		JudgedCache[hashKey] = true
		subRet, err := processFuncDecl(path, &subFuncDecl, rule)
		if err != nil {
			return false, err
		} else if subRet {
			return true, nil
		} else {
			JudgedCache[hashKey] = false
			return false, nil
		}
	} else {
		JudgedCache[hashKey] = false
		return false, nil
	}
	return true, nil
}

// processNameDecl 返回true 表示没有鉴权框架
func processNameDecl(name string, rule *rule.Rule, path string) (bool, error) {
	if ret, ok := nameJudgedCache[name]; ok && ret {
		return true, nil
	} else if ok {
		// 包含了鉴权框架
		return false, nil
	}
	if units, ok := CodeCache[name]; ok {
		for _, unit := range units {
			hashValue := GetFuncAstHash(unit.FuncAst)
			if ret, ok := JudgedCache[hashValue]; ok {
				if ret {
					continue
				}
				return false, nil
			}

			if ret, err := matchStr(rule.GoModeTargetRule.Rule, string(unit.Code)); err != nil {
				logger.Error(err.Error())
			} else if ret { // 发现该层没有目标的函数，将会向下查询调用
				hashValue = GetFuncAstHash(unit.FuncAst)
				JudgedCache[hashValue] = true
				//startLine, endLine := getStartAndEndLine(api.Fset, api.FuncAst)
				//result := fmt.Sprintf("%d:%d:%s", startLine, endLine, api.FuncAst.Name.Name)
				//idors = append(idors, result)
				subRet, err := processFuncDecl(path, unit.FuncAst, rule)
				if err != nil {
					logger.Error(err.Error())
					continue
				} else if subRet {
					continue
				} else {
					JudgedCache[hashValue] = false
					return false, nil
				}
			} else {
				JudgedCache[hashValue] = false
				return false, nil
			}
		}
		nameJudgedCache[name] = true
	} else {
		nameJudgedCache[name] = true
	}
	return true, nil
}

func scanDecls(asts []ast.Decl, f func(decl *ast.FuncDecl)) {
	for _, decl := range asts {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			f(funcDecl)
		}
	}
}
