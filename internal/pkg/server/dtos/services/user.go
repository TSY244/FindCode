package services

type AddUserServiceParam struct {
	Password string `json:"password"`
	Role     int    `json:"role"`
	UserName string `json:"username"`
}
