package users

import (
	"cogmoteHub/internal/models"

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
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) GetByUID(uid uint64) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, uid).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) GetByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) List() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
