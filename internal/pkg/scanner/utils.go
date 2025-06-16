package scanner

import (
	"ScanIDOR/pkg/color"
	"ScanIDOR/pkg/logger"
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
func printResult(ret map[string][]string) {
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
	for filePath, ret := range Result {
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
