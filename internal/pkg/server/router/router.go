package router

import (
	"ScanIDOR/internal/pkg/server/controllers"
	"ScanIDOR/internal/pkg/server/services"
	"github.com/gin-gonic/gin"
)

func SetupRouter(findCodeService services.FindCodeService, user services.UserService) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// index
	r.GET("/", controllers.NewRootController().Root)

	// findCode group
	findCodeController := controllers.NewFindCodeController(findCodeService)
	findCode := r.Group("/findCode")
	{
		findCode.POST("/scan", findCodeController.Scan)
		findCode.GET("/tool", findCodeController.ShowScanHtml)
	}
	// user group
	userController := controllers.NewUserController(user)
	userGroup := r.Group("/user")
	{
		userGroup.POST("/add", userController.AddUser)
		userGroup.GET("/create", userController.ShowCreateUserHtml)
	}
	return r
}
