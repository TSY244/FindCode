package fingerprint

import (
	"fmt"
)

func GetProductPrint(productPath string) (string, error) {
	// 遍历寻找 go.mod
	fileName := map[string]struct{}{
		"go.mod": {},
	}
	filePath := getFilePath(productPath, fileName)
	goModeFilePath, ok := filePath["go.mod"]
	if !ok {
		return "", fmt.Errorf("go.mod not found")
	}
	ret := getAllPrints(goModeFilePath)
	if len(ret) == 1 {
		return ret[0], nil
	}
	// 只要满足就直接返回
	param := detailFuncParam{
		ProductPathKey: productPath,
	}
	for _, v := range ret {
		f, ok := detailFuncMap[v]
		if !ok {
			continue
		}
		if bRet := f(param); bRet {
			return v, nil
		}
	}
	return "", fmt.Errorf("unsupported frameworks")
}
