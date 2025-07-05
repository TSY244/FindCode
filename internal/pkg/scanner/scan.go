package scanner

import (
	"ScanIDOR/internal/pkg/rule"
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/pkg/logger"
	"ScanIDOR/utils/util"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

type modeFunc func(path string, target *rule.Rule, env *Env) error

var (
	ModeFuncMap = map[string]modeFunc{
		consts.GoMode:  goModeFunc,
		consts.StrMode: strModeFunc,
	}
	strModeFlag  = false
	isUsedGoMode = false
)

// Scan
//
// 参数：
//   - fcFlag 用于兼容旧版本，但是需要改造成后端项目做的取舍
func Scan(path string, r *rule.Rule, env *Env) error {
	if r == nil {
		return errors.New("r is nil")
	}
	//LoadCtx(ctx, r)

	info, err := checkFileStatue(path)
	if err != nil {
		return err
	}
	// 根据文件类型启动不同的扫描逻辑
	if info.IsDir() {
		if err := dealDir(path, r, consts.BeginLevel, env); err != nil {
			return err
		}
	} else {
		if err := dealFile(path, r, env); err != nil {
			return err
		}
	}
	if isUseMode(r, consts.GoMode) {
		isUsedGoMode = true
		if err := processFuncDecls(r, env); err != nil {
			return err
		}
	}

	// 判断是否使用ai mode
	if isUseMode(r, consts.AiMode) {
		if r.GoModeTargetRule.Rule != "true" {
			r.AiConfig.Prompt = r.AiConfig.Prompt + "\n\n相关鉴权函数逻辑" + r.GoModeTargetRule.Rule
		}
		if err := aiScan(r.AiConfig, env); err != nil {
			return err
		}
	}

	SaveToFile(r.TaskName, env)
	printResult(env.Result, env)
	logger.Info("模块已经扫描结束")
	return nil
}

func isUseMode(rule *rule.Rule, target string) bool {
	for _, r := range rule.Mode {
		if r == target {
			return true
		}
	}
	return false
}

// dealDir 将目录下的文件逐个遍历，将其中的函数保存到对应的cache 中
func dealDir(path string, target *rule.Rule, deep int, env *Env) error {
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
			if err := dealDir(filepath.Join(path, file.Name()), target, deep+1, env); err != nil {
				return err
			}
		} else {
			filePath := filepath.Join(path, file.Name())

			if err := dealFile(filePath, target, env); err != nil {
				return err
			}
		}
	}
	return nil
}

// dealFile 解析文件，提取出函数，保存到cache 中
func dealFile(path string, target *rule.Rule, env *Env) error {
	info, err := checkFileStatue(path)
	if err != nil {
		return err
	}
	if !isGoFile(info) {
		return nil
	}

	for _, mode := range target.Mode {
		if mode == consts.AiMode {
			continue
		}
		f, ok := ModeFuncMap[mode]
		if !ok {
			logger.Warn("不存在mode: ", mode)
			continue
		}
		if err := f(path, target, env); err != nil {
			return err
		}
	}
	return nil
}

func strModeFunc(path string, r *rule.Rule, env *Env) error {
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

func goModeFunc(path string, target *rule.Rule, env *Env) error {
	err := processFile(path, target, env)
	if err != nil {
		return err
	}
	return nil
}

// processFile 解析文件真正的函数
func processFile(filePath string, rule *rule.Rule, env *Env) error {
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

	//func 定义了对 解析的函数节点 的处理逻辑
	scanDecls(fileAst.Decls, func(decl *ast.FuncDecl) {
		// 处理的是，将api func 和 非api func 补充到chache
		ret, err := filter(decl, rule)

		if !errors.Is(err, ArgsSizeNotEqualErr) && !errors.Is(err, FuncParamsNotEqualErr) && err != nil {
			return
		}
		// path 是否存在规则, 如果有则对path 进行判断，
		if rule.Path.Rule != "" && rule.Path.Rule != "true" {
			if subRet, err := matchStr(rule.Path.Rule, filePath); err != nil {
				return
			} else if !subRet {
				ret = false
			}
		}

		if ret {
			if err := saveApiFunc(filePath, srcStr, decl, fset, env); err != nil {
				logger.Error(err.Error())
				return
			}
		} else {
			savaNoApiFunc(filePath, decl, srcStr, fset, env)
		}
	})
	return nil
}

func processFuncDecls(rule *rule.Rule, env *Env) error {
	for path, apis := range env.ApiCache {
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
			ret, err := processApi(api, rule, path, env)
			if err != nil {
				return err
			}
			result = append(result, ret...)
		}
		if len(result) > 0 {
			env.Result[path] = result
		}
	}
	return nil
}

func processApi(api cacheUnit, rule *rule.Rule, path string, env *Env) ([]string, error) {
	//debug
	idors := make([]string, 0)
	funcCode, err := util.Decompress(api.Code)
	if err != nil {
		return nil, err
	}

	if ret, err := matchStr(rule.GoModeTargetRule.Rule, string(funcCode)); err != nil {
		logger.Error(err.Error())
	} else if ret { // 处理子调用逻辑
		subRet, err := processFuncDecl(path, api.FuncAst, rule, env)
		if err != nil {
			return nil, err
		}
		if subRet {
			// 可能存在问题的，统一保存
			hashValue := GetFuncAstHash(api.FuncAst)
			env.JudgedCache[hashValue] = true
			startLine, endLine := getStartAndEndLine(api.Fset, api.FuncAst)
			result := fmt.Sprintf("%d:%d:%s", startLine, endLine, api.FuncAst.Name.Name)
			idors = append(idors, result)
		}
	}
	return idors, nil
}

// processFuncDecl 返回true 表示没有鉴权框架
func processFuncDecl(path string, decl *ast.FuncDecl, rule *rule.Rule, env *Env) (bool, error) {
	// 1. 获取所有的子调用的函数名字
	allSubFuncDecls, _ := getAllSubFuncDecls(decl, path, env)

	// 2. 无法获取funcdecl 的子调用，通过name 进行判断
	// todo： 和funcdecl 分开统计
	//for _, name := range names {
	//	if ret, err := processNameDecl(name, rule, path, env); err != nil {
	//		return false, err
	//	} else if !ret {
	//		return false, nil
	//	}
	//}

	// 统计funcdecl 的调用逻辑
	for _, subFuncDecl := range allSubFuncDecls {
		if ret, err := processSubFuncDecl(subFuncDecl, rule, path, env); err != nil {
			return false, err
		} else if !ret {
			return false, nil
		}
	}
	return true, nil
}

// processSubFuncDecl 处理函数调用的子函数的
func processSubFuncDecl(subFuncDecl ast.FuncDecl, rule *rule.Rule, path string, env *Env) (bool, error) {
	hashKey := GetFuncAstHash(&subFuncDecl)
	if cacheRet, ok := env.JudgedCache[hashKey]; ok {
		if cacheRet {
			return true, nil
		}
		return cacheRet, nil
	}

	unit, ok := env.FuncCacheMap[hashKey]
	if !ok {
		if unit == nil {
			return false, nil
		}
		name := unit.FuncAst.Name.Name
		if ret, err := processNameDecl(name, rule, path, env); err != nil {
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
		env.JudgedCache[hashKey] = true
		subRet, err := processFuncDecl(path, &subFuncDecl, rule, env)
		if err != nil {
			return false, err
		} else if subRet {
			return true, nil
		} else {
			env.JudgedCache[hashKey] = false
			return false, nil
		}
	} else {
		env.JudgedCache[hashKey] = false
		return false, nil
	}
	return true, nil
}

// processNameDecl 返回true 表示没有鉴权框架
func processNameDecl(name string, rule *rule.Rule, path string, env *Env) (bool, error) {
	if ret, ok := env.nameJudgedCache[name]; ok && ret {
		return true, nil
	} else if ok {
		// 包含了鉴权框架
		return false, nil
	}
	if units, ok := env.CodeCache[name]; ok {
		for _, unit := range units {

			// ------------
			// 测试code Name 代码拿funcdelc
			//subCode, names := getAllSubFuncDecls(unit.FuncAst, unit.FilePath)
			//fmt.Println(subCode, names)
			// ------------

			hashValue := GetFuncAstHash(unit.FuncAst)
			if ret, ok := env.JudgedCache[hashValue]; ok {
				if ret {
					continue
				}
				return false, nil
			}

			if ret, err := matchStr(rule.GoModeTargetRule.Rule, string(unit.Code)); err != nil {
				logger.Error(err.Error())
			} else if ret { // 发现该层没有目标的函数，将会向下查询调用
				hashValue = GetFuncAstHash(unit.FuncAst)
				env.JudgedCache[hashValue] = true
				//startLine, endLine := getStartAndEndLine(controllers.Fset, controllers.FuncAst)
				//result := fmt.Sprintf("%d:%d:%s", startLine, endLine, controllers.FuncAst.Name.Name)
				//idors = append(idors, result)
				subRet, err := processFuncDecl(path, unit.FuncAst, rule, env)
				if err != nil {
					logger.Error(err.Error())
					continue
				} else if subRet {
					continue
				} else {
					env.JudgedCache[hashValue] = false
					return false, nil
				}
			} else {
				env.JudgedCache[hashValue] = false
				return false, nil
			}
		}
		env.nameJudgedCache[name] = true
	} else {
		env.nameJudgedCache[name] = true
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

// 获取当前func 底层的被调用者的函数
//func getAllSubCode(decl *ast.FuncDecl, path string, env *Env) ([]string, error) {
//	subDecl, names := getAllSubFuncDecls(decl, path)
//	var allSubCode []string
//	for _, sub := range subDecl {
//		funcCode := getFuncCode(path, &sub)
//		if funcCode == "" {
//			continue
//		}
//		allSubCode = append(allSubCode, funcCode)
//	}
//	for _, name := range names {
//		if units, ok := env.CodeCache[name]; ok {
//			for _, unit := range units {
//				allSubCode = append(allSubCode, string(unit.Code))
//			}
//		}
//	}
//	return allSubCode, nil
//}

func getAllSubCodeWithLevel(decl *ast.FuncDecl, path string, env *Env, level, maxLevel uint) ([]string, error) {
	if level == 0 {
		return nil, errors.New("level is zero")
	}
	if level > maxLevel {
		return nil, nil
	}
	subDecl, names := getAllSubFuncDecls(decl, path, env)
	var allSubCode []string
	for _, sub := range subDecl {
		hashKey := GetFuncAstHash(&sub)
		funcInfo, ok := env.FuncCacheMap[hashKey]
		if !ok {
			continue
		}
		//funcCode := getFuncCode(funcInfo.FilePath, &sub)
		//if funcCode == "" {
		//	continue
		//}
		//funcCodeEncrypted := funcInfo.Code
		//funcCodeBytes, err := util.Decompress(funcCodeEncrypted)
		//if err != nil {
		//	logger.Error(err.Error())
		//	continue
		//}
		//funcCode := string(funcCodeBytes)
		funcCode := string(funcInfo.Code)
		if level != consts.FirstLevel {
			funcCode += decl.Name.Name + " 调用的代码如下: " + funcCode
		}
		allSubCode = append(allSubCode, funcCode)
		// 添加底层的代码
		nextLevelSubCode, err := getAllSubCodeWithLevel(funcInfo.FuncAst, path, env, level+1, maxLevel)
		if err != nil {
			return nil, err
		}
		allSubCode = append(allSubCode, nextLevelSubCode...)
	}
	for _, name := range names {
		if units, ok := env.CodeCache[name]; ok {
			for _, unit := range units {
				allSubCode = append(allSubCode, string(unit.Code))
			}
		}
	}
	return allSubCode, nil
}
