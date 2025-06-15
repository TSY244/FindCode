package repositories

import (
	"ScanIDOR/internal/pkg/server/models"
	"context"
	"gorm.io/gorm"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, token *models.Token) error
}

type TokenRepoImpl struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) TokenRepository {
	return &TokenRepoImpl{
		db: db,
	}
}

func (repo *TokenRepoImpl) CreateToken(ctx context.Context, token *models.Token) error {
	return repo.db.WithContext(ctx).Select("Token").Create(token).Error
}
