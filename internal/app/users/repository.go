package users

import (
	"cogmoteHub/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByUID(uid uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	List() ([]models.User, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db}
}

func (r *repo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repo) GetByID(id uuid.UUID) (*models.User, error) {
	ctx := context.Background()

	user, err := gorm.G[models.User](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) GetByUID(uid uint64) (*models.User, error) {
	ctx := context.Background()

	user, err := gorm.G[models.User](r.db).Where("uid = ?", uid).First(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) GetByEmail(email string) (*models.User, error) {
	ctx := context.Background()

	user, err := gorm.G[models.User](r.db).Where("email = ?", email).First(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) List() ([]models.User, error) {
	ctx := context.Background()

	users, err := gorm.G[models.User](r.db).Find(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
