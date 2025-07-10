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

var (
	PermissionModel = map[string]string{
		consts.PublicPermission:              consts.PublicPermissionDetail,
		consts.UserGroupPermission:           consts.UserGroupPermissionDetail,
		consts.RBACPermission:                consts.RBACPermissionDetail,
		consts.RBACWithConstraintsPermission: consts.RBACWithConstraintsPermissionDetail,
		consts.ABACOrPBACPermission:          consts.ABACOrPBACPermissionDetail,
		consts.DACPermission:                 consts.DACPermissionDetail,
		consts.MACPermission:                 consts.MACPermissionDetail,
	}
)
