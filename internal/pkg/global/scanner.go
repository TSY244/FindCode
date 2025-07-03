package global

import (
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/pkg/fingerprint"
)

var (
	RuleMap = map[string]string{
		fingerprint.GinPrint:  consts.GinRule,
		fingerprint.GoSwagger: consts.GoSwaggerRule,
		fingerprint.TRPCPrint: consts.TrpcRule,
	}
)
