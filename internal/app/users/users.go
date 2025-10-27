package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r gin.IRouter, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	handler.RegisterRoutes(r)
}
