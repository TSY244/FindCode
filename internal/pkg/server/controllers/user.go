package controllers

import (
	"ScanIDOR/internal/pkg/server/dtos/requests"
	serviceParam "ScanIDOR/internal/pkg/server/dtos/services"
	"ScanIDOR/internal/pkg/server/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (u *UserController) AddUser(c *gin.Context) {
	req := requests.AddUserRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg.go":   "参数错误",
			"go.error": "",
		})
		return
	}
	param := serviceParam.AddUserServiceParam{
		UserName: req.UserName,
		Password: req.Password,
		Role:     req.Role,
	}
	if err := u.userService.AddUser(c, &param); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"msg.go":   "添加用户失败",
			"go.error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "success.html", gin.H{
		"msg.go": "添加用户成功",
	})
}

func (u *UserController) ShowCreateUserHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "create_user.html", nil)
}
