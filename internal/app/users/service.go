package users

import (
	"cogmoteHub/internal/encrypt"
	"cogmoteHub/internal/models"
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	Register(req RegisterRequest) (*UserResponse, error)
	GetAll() ([]UserResponse, error)
	GetByID(id uuid.UUID) (*UserResponse, error)
	GetByUID(uid uint64) (*UserResponse, error)
	GetByEmail(email string) (*UserResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Register(req RegisterRequest) (*UserResponse, error) {
	if _, err := s.repo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	hashed, err := encrypt.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashed,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		UID:      user.UID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *service) GetAll() ([]UserResponse, error) {
	var ras []UserResponse
	users, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		ras = append(ras, UserResponse{
			ID:       user.ID,
			UID:      user.UID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return ras, nil
}

func (s *service) GetByID(id uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		UID:      user.UID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *service) GetByUID(uid uint64) (*UserResponse, error) {
	user, err := s.repo.GetByUID(uid)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		UID:      user.UID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *service) GetByEmail(email string) (*UserResponse, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		UID:      user.UID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
