package fingerprint

type detailFuncParam map[string]interface{}

type detailFunc func(param detailFuncParam) bool

var (
	detailFuncMap = map[string]detailFunc{
		TRPCPrint: trpcDetailFunc,
		GoSwagger: goSwaggerDetailFunc,
		GinPrint:  ginDetailFunc,
	}
)

func trpcDetailFunc(param detailFuncParam) bool {
	// 规则之间是or 规则内是and
	// 1. 项目中必须存在
	fileName := "trpc_go.yaml"
	fileNames := map[string]struct{}{
		fileName: {},
	}
	filePath := getFilePath(param[ProductPathKey].(string), fileNames)
	if filePath != nil {
		_, ok := filePath[fileName]
		if ok {
			return true
		}
	}
	return false
}

func goSwaggerDetailFunc(param detailFuncParam) bool {
	fileName := "swagger.yaml"
	fileNames := map[string]struct{}{
		fileName: {},
	}
	filePath := getFilePath(param[ProductPathKey].(string), fileNames)
	if filePath != nil {
		_, ok := filePath[fileName]
		if ok {
			return true
		}
	}
	return false
}

func ginDetailFunc(param detailFuncParam) bool {
	targetContent := "*gin.Context"
	return IsContainContent(param[ProductPathKey].(string), targetContent)
}
