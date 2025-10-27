package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r gin.IRouter) {
	auth := r.Group("/auth")
	// TODO: Passkeys
	// auth.GET("/register/challenge")
	// auth.POST("/register")
	// auth.GET("/login/challenge")
	// auth.POST("/auth/login")
	// auth.GET("/passkeys")
	// auth.DELETE("/passkeys/credentialId")

	auth.POST("refresh")
	auth.POST("logout")
	auth.GET("me")

	password := auth.Group("/password")
	password.POST("/register")
	password.POST("/login")
}
