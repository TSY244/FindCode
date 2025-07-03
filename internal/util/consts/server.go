package consts

import "ScanIDOR/pkg/fingerprint"

const (
	ruleDir       = "rule/"
	ginApiYaml    = "find_gin_api.yaml"
	goSwaggerYaml = "find_go_swagger_api.yaml"
	trpcYaml      = "find_trpc_api.yaml"
)

var (
	TypeMap = map[string]string{
		fingerprint.GinPrint:  ruleDir + ginApiYaml,
		fingerprint.GoSwagger: ruleDir + goSwaggerYaml,
		fingerprint.TRPCPrint: ruleDir + trpcYaml,
	}
)
