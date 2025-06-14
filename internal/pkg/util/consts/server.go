package consts

const (
	ruleDir       = "rule/"
	ginApiYaml    = "find_gin_api.yaml"
	goSwaggerYaml = "find_go_swagger_api.yaml"
	trpcYaml      = "find_trpc_api.yaml"
)

var (
	TypeMap = map[string]string{
		"gin":        ruleDir + ginApiYaml,
		"go_swagger": ruleDir + goSwaggerYaml,
		"trpc":       ruleDir + trpcYaml,
	}
)
