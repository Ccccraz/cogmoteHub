package login

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.POST("/login", h.PostLogin)
	r.POST("/refresh", h.PostRefresh)
}

func (h *Handler) PostLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	rep := LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	c.JSON(http.StatusCreated, rep)
}

func (h *Handler) PostRefresh(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	refreshToken := parts[1]

	tokenPair, err := h.service.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	rep := LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	c.JSON(http.StatusCreated, rep)
}
