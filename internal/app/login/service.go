package login

import (
	"cogmoteHub/internal/authenticator"
	"cogmoteHub/internal/encrypt"
	"cogmoteHub/internal/models"
	"errors"
	"time"
)

type Service struct {
	repo      Repo
	tokenAuth authenticator.JwtAuthenticator
}

func NewService(repo Repo, tokenAuth authenticator.JwtAuthenticator) *Service {
	return &Service{
		repo:      repo,
		tokenAuth: tokenAuth,
	}
}

func (s *Service) Login(req LoginRequest) (*authenticator.TokenPair, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	ok, err := encrypt.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("invalid credentials")
	}

	tokenPair, err := s.tokenAuth.IssueTokens(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken := models.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: tokenPair.RefreshTokenExpiresAt,
	}

	s.repo.CreateRefreshToken(&refreshToken)

	return tokenPair, nil
}

func (s *Service) Refresh(refreshToken string) (*authenticator.TokenPair, error) {
	storedRefreshToken, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	err = s.repo.DeleteRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if storedRefreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	tokenPair, err := s.tokenAuth.IssueTokens(storedRefreshToken.UserID)
	if err != nil {
		return nil, err
	}

	newRefreshToken := models.RefreshToken{
		UserID:    storedRefreshToken.UserID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: tokenPair.RefreshTokenExpiresAt,
	}

	s.repo.CreateRefreshToken(&newRefreshToken)

	return tokenPair, nil
}
