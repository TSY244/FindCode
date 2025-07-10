package utils

import (
	"ScanIDOR/internal/pkg/global"
	"ScanIDOR/internal/util/consts"
	"fmt"
)

func GetProductTypeDetail(productType string) string {
	return consts.ProductDetail + productType
}

func GerPermissionDetail(permissionModel string) string {
	detail, ok := global.PermissionModel[permissionModel]
	if !ok {
		return ""
	}
	return fmt.Sprintf(consts.PermissionDetail, permissionModel, detail)

}
