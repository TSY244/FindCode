package repositories

import (
	"ScanIDOR/internal/pkg/server/models"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByToken(ctx context.Context, token string) (*models.User, error)
}

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &UserRepoImpl{
		db: db,
	}
}

func (repo *UserRepoImpl) CreateUser(ctx context.Context, user *models.User) error {
	return repo.db.WithContext(ctx).Create(user).Error
}

func (repo *UserRepoImpl) FindUserByToken(ctx context.Context, token string) (*models.User, error) {
	query := repo.db.WithContext(ctx).Preload("token")
	var user models.User
	if err := query.Where("token.token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
