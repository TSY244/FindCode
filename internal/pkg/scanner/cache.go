package scanner

import (
	"ScanIDOR/pkg/logger"
	"ScanIDOR/utils/util"
	"go/ast"
	"go/token"
)

type cacheUnit struct {
	FuncAst  *ast.FuncDecl
	Fset     *token.FileSet
	Code     []byte
	FilePath string
}

type Cache map[string][]cacheUnit

// savaNoApiFunc 保存所有的非api 函数
func savaNoApiFunc(filePath string, decl *ast.FuncDecl, srcStr string, fset *token.FileSet, env *Env) {
	hashKey := GetFuncAstHash(decl)
	code := getFuncCode(srcStr, decl)
	unit := cacheUnit{
		FuncAst:  decl,
		Fset:     fset,
		Code:     []byte(code),
		FilePath: filePath,
	}
	env.FuncCacheMap[hashKey] = &unit
	env.CodeCache[decl.Name.Name] = append(env.CodeCache[decl.Name.Name], &unit)
}

// saveApiFunc 保存api 函数
func saveApiFunc(filePath, srcStr string,
	decl *ast.FuncDecl, fset *token.FileSet, env *Env) error {
	code, err := util.Compress([]byte(getFuncCode(srcStr, decl)))
	if err != nil {
		logger.Errorf("在处理api 函数的时候，保存错误")
		return err
	}

	apis := env.ApiCache[filePath]
	apis = append(apis, cacheUnit{
		FuncAst: decl,
		Code:    code,
		Fset:    fset,
	})
	env.ApiCache[filePath] = apis
	return nil
}
