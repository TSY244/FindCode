package requests

type AddUserRequest struct {
	UserName string `form:"userName" binding:"required,min=3,max=50"`
	Password string `form:"password" binding:"required,min=6,max=100"`
	Role     int    `form:"role" binding:"omitempty,min=1,max=2"` // 1-普通用户 2-管理员
}
