package services

import (
	serviceParam "ScanIDOR/internal/pkg/server/dtos/services"
	"ScanIDOR/internal/pkg/server/models"
	"ScanIDOR/internal/pkg/server/repositories"
	"ScanIDOR/utils/util"
	"context"
)

type UserService interface {
	AddUser(ctx context.Context, param *serviceParam.AddUserServiceParam) error
	FindUserByToken(ctx context.Context, token string) (*models.User, error)
}

type UserServiceImpl struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.TokenRepository
}

func NewUserService(userRepo repositories.UserRepository, token repositories.TokenRepository) UserService {
	return &UserServiceImpl{
		userRepo:  userRepo,
		tokenRepo: token,
	}
}

func (s *UserServiceImpl) AddUser(ctx context.Context, param *serviceParam.AddUserServiceParam) error {
	var user models.User
	var token models.Token
	var err error
	user.UserName = param.UserName
	user.Password, err = util.HashPassword(param.Password)
	if err != nil {
		return err
	}
	user.Role = param.Role
	token.Token = util.GenerateToken()
	user.UserToken = token
	return s.userRepo.CreateUser(ctx, &user)
	//
	//err = s.userRepo.CreateUser(ctx, &user)
	//if err != nil {
	//	return err
	//}
	//
	//return s.tokenRepo.CreateToken(ctx, &token)
}

func (s *UserServiceImpl) FindUserByToken(ctx context.Context, token string) (*models.User, error) {
	return s.userRepo.FindUserByToken(ctx, token)
}
