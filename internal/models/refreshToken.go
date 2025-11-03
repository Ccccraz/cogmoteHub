package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Base
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	UserID    uuid.UUID
}
