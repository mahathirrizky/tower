package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/user/tower-tracker-bima/backend/controllers"
	"github.com/user/tower-tracker-bima/backend/middleware" // Import middleware
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		// Apply rate limiting to login and change-password
		auth.POST("/login", middleware.RateLimitMiddleware(5, 1), controllers.Login)
		auth.POST("/register", controllers.Register)

		// Authenticated routes
		auth.Use(middleware.AuthMiddleware()) // Apply auth middleware to subsequent routes in this group
		auth.PUT("/change-password", middleware.RateLimitMiddleware(5, 1), controllers.ChangePassword) // New route for changing password
	}
}
