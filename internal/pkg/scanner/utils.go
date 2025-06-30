package scanner

import (
	"ScanIDOR/internal/config"
	"ScanIDOR/internal/pkg/ai"
	"ScanIDOR/internal/pkg/rule"
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/pkg/color"
	"ScanIDOR/pkg/logger"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
)

// GetFuncAstHash 通过func 性质生成hash
func GetFuncAstHash(funcAst *ast.FuncDecl) string {
	name := funcAst.Name.Name
	paramNames, _ := getFuncParamNames(funcAst)
	paramTypes, _ := getFuncParamTypes(funcAst)
	rawSplice := append(paramNames, paramTypes...)
	rawSplice = append(rawSplice, name)
	rawStr := strings.Join(rawSplice, "&")
	// 求字符串的hash值
	hasher := sha256.New()
	hasher.Write([]byte(rawStr))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

// checkFileStatue 检查文件是否存在
func checkFileStatue(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Errorf("file %s does not exist", path)
		}
		return nil, err
	}
	return info, nil
}

func IsEqual(raw *ast.FuncDecl, target *ast.FuncDecl) (bool, error) {
	if raw.Name.Name != target.Name.Name {
		return false, nil
	}
	rawNams, err := getFuncParamNames(raw)
	if err != nil {
		return false, err
	}
	targetNams, err := getFuncParamNames(target)
	if err != nil {
		return false, err
	}
	if len(rawNams) != len(targetNams) {
		return false, FuncParamsNotEqualErr
	}

	rawTypes, err := getFuncParamTypes(raw)
	if err != nil {
		return false, err
	}
	targetTypes, err := getFuncParamTypes(target)
	if err != nil {
		return false, err
	}
	if len(rawTypes) != len(targetTypes) {
		return false, FuncParamsNotEqualErr
	}

	return true, nil
}

// readFile2String 读取文件内容
func readFile2String(path string) (string, error) {
	fileP, err := os.Open(path)
	if err != nil {
		return "", err
	}
	fileBytes, err := io.ReadAll(fileP)
	if err != nil {
		return "", err
	}
	srcStr := string(fileBytes)
	return srcStr, nil
}

func isGoFile(f os.FileInfo) bool {
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

// 获取目录下所有的文件
func getAllFilesInDir(path string) ([]os.FileInfo, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func getAllSubFuncDecls(decl *ast.FuncDecl, path string) (map[string]ast.FuncDecl, []string) {
	funcDecls := make(map[string]ast.FuncDecl)
	nameSet := make(map[string]struct{})
	names := make([]string, 0)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil
	}

	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if decl.Name.Name != fn.Name.Name {
			return true
		}
		// 在函数体内查找所有函数调用
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			callName := getCallName(call.Fun)
			if callName == "" {
				//fmt.Println("调用函数:", callName)
				return true
			}
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Obj == nil {
				nameSet[callName] = struct{}{}
				return true
			}

			if funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl); ok {
				hashKey := GetFuncAstHash(ident.Obj.Decl.(*ast.FuncDecl))
				funcDecls[hashKey] = *funcDecl
			}
			return true
		})
		return false
	})
	for k := range nameSet {
		//if k == decl.Name.Name || k == "" {
		//	continue
		//}
		if k == "" {
			continue
		}
		names = append(names, k)
	}

	return funcDecls, names
}

func getCallName(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return v.Sel.Name
	default:
		return ""
	}
}

func getStartAndEndLine(fset *token.FileSet, desl *ast.FuncDecl) (int, int) {
	startLine := fset.Position(desl.Pos()).Line
	endLine := fset.Position(desl.End()).Line
	return startLine, endLine
}

// 获取func 的代码
func getFuncCode(srcStr string, decl *ast.FuncDecl) string {
	start := decl.Pos()
	end := decl.End()
	var funcCode string
	if len(srcStr) < int(start) {
		return ""
	}
	if int(end) < len(srcStr) {
		funcCode = srcStr[start:end]
	} else {
		funcCode = srcStr[start:]
	}
	return funcCode
}

// printResult 打印结果
func printResult(ret map[string][]string, env *Env) {
	if strModeFlag {
		color.Magenta("该项目存在filter\n\n")
	} else {
		color.Magenta("该项目不存在filter\n\n")
	}

	allSize := 0

	if !isUsedGoMode {
		return
	}
	if len(ret) == 0 {
		logger.Infoln("go mode没有扫描出结果！")
		return
	}
	for filePath, ret := range env.Result {
		//logger.Warn(filePath + " 的文件可能存在的越权漏洞如下：")
		color.HRed(filePath + " 的文件可能存在的越权漏洞如下：\n")
		allSize += len(ret)
		for _, str := range ret {
			//logger.Warn(str)
			color.Yellow(str + "\n")
		}
		fmt.Println()
	}
	logger.Infof("函数个数: %d", allSize)
}

// LoadCtx 暂时只负责处理ai scan 的传参问题
func LoadCtx(ctx context.Context, r *rule.Rule) {
	isUseCtx, ok := ctx.Value(consts.IsUseCtxKey).(bool)
	if !ok {
		return
	}
	if !isUseCtx { // 当不使用ctx 说明是cli 模式。
		if config.CoreConfig == nil || config.CoreConfig.AiConfig == nil {
			return
		}
		r.AiConfig = config.CoreConfig.AiConfig
		return
	}
	aiConFig, ok := ctx.Value(consts.AiConfigKey).(*ai.Config)
	if ok {
		r.AiConfig = aiConFig
	}

}

func GetResult(clonePath string, env *Env) (*map[string][]string, *map[string]AiBoolResultUnitWithStatue) {
	if strings.HasPrefix(clonePath, "./") {
		clonePath = strings.Replace(clonePath, "./", "", 1)
	}
	locker.RLock()
	defer locker.RUnlock()
	ret := make(map[string][]string)
	boolRet := make(map[string]AiBoolResultUnitWithStatue)
	for k, v := range env.Result {
		if strings.HasPrefix(k, clonePath) {
			ret[k] = v
		}
	}
	for k, v := range env.AiBoolResult {
		if strings.HasPrefix(k, clonePath) {
			boolRet[k] = getAiBoolResultUnitWithStatue(&v)
		}
	}
	return &ret, &boolRet
}

func getAiBoolResultUnitWithStatue(aiBoolResultUnit *AiBoolResultUnit) AiBoolResultUnitWithStatue {
	aiBoolResult := make(AiBoolResultUnitWithStatue)
	for k, v := range *aiBoolResultUnit {
		aiBoolResult[k] = aiBoolUnitWithStatue{
			Statue:      getStatue(&v),
			AiBoolUnits: v,
		}
	}
	return aiBoolResult
}

func getStatue(abus *[]aiBoolUnit) int {
	var statue int
	size := len(*abus)
	for _, unit := range *abus {
		if unit.Result == "true" {
			statue++
		} else if unit.Result == "false" {
			statue--
		}
	}
	if statue == size {
		return 1
	} else if statue == -size {
		return -1
	}
	return 0
}

func ClearResult(clonePath string, env *Env) {
	if strings.HasPrefix(clonePath, "./") {
		clonePath = strings.Replace(clonePath, "./", "", 1)
	}
	locker.Lock()
	defer locker.Unlock()
	for k := range env.Result {
		if strings.HasPrefix(k, clonePath) {
			delete(env.Result, k)
		}
	}
	for k := range env.AiBoolResult {
		if strings.HasPrefix(k, clonePath) {
			delete(env.AiBoolResult, k)
		}
	}
}

func GetContainsRule(strs []string) string {
	var allFuncs []string
	for _, str := range strs {
		allFuncs = append(allFuncs, "contain("+str+")")
	}
	ruleStr := strings.Join(allFuncs, " || ")
	return "! " + ruleStr
}
