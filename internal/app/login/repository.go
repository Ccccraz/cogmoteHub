package login

import (
	"cogmoteHub/internal/models"
	"context"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetByEmail(email string) (*models.User, error) {
	ctx := context.Background()

	user, err := gorm.G[models.User](r.db).Where("email = ?", email).First(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) CreateRefreshToken(refreshToken *models.RefreshToken) error {
	ctx := context.Background()

	return gorm.G[models.RefreshToken](r.db).Create(ctx, refreshToken)
}

func (r *Repo) GetRefreshToken(refreshToken string) (*models.RefreshToken, error) {
	ctx := context.Background()

	token, err := gorm.G[models.RefreshToken](r.db).Where("token = ?", refreshToken).First(ctx)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *Repo) DeleteRefreshToken(refreshToken string) error {
	ctx := context.Background()

	_, err := gorm.G[models.RefreshToken](r.db).Where("token = ?", refreshToken).Delete(ctx)

	return err
}
