package users

import "github.com/google/uuid"

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	UID      uint64    `json:"uid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}
