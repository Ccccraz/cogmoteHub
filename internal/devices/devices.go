package devices

import (
	"errors"
	"net/http"

	"cogmoteHub/internal/db"
	"cogmoteHub/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetDevicesHandler(c *gin.Context) {
	var devices []*models.Device
	database := db.Get()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not initialized"})
		return
	}

	if err := database.Preload("Animals").Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, devices)
}

func GetDeivceHandler(c *gin.Context) {
	id := c.Param("id")

	var device models.Device
	database := db.Get()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not initialized"})
		return
	}

	if err := database.Preload("Animals").First(&device, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

func PostDeviceHandler(c *gin.Context) {
	var payload models.Device
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database := db.Get()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not initialized"})
		return
	}

	payload.ID, _ = uuid.NewV7()

	if err := database.Create(&payload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payload)
}

func RegisterRoutes(r gin.IRouter) {
	r.GET("/devices", GetDevicesHandler)
	r.POST("/devices", PostDeviceHandler)
	r.GET("/devices/:id", GetDeivceHandler)
}
