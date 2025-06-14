package router

import (
	"ScanIDOR/internal/pkg/server/api"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// exec findCode
	r.GET("/", api.Root)
	r.POST("/scan", api.Scan)
}
