package models

type User struct {
	Base
	UID          uint64 `gorm:"not null;autoIncrement;uniqueIndex"`
	Username     string `gorm:"not null;unique"` // Username
	Email        string `gorm:"uniqueIndex"`     // Email
	PasswordHash string `gorm:"not null"`        // Password hash
	RefreshToken []RefreshToken
}
