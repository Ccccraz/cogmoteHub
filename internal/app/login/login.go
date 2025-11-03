package login

import (
	"cogmoteHub/internal/authenticator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r gin.IRouter, db *gorm.DB, tokenAuth *authenticator.JwtAuthenticator) {
	repo := NewRepository(db)
	service := NewService(*repo, *tokenAuth)
	handler := NewHandler(*service)
	handler.RegisterRoutes(r)
}
