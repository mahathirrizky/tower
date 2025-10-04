package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/user/tower-tracker-bima/backend/controllers"
	"github.com/user/tower-tracker-bima/backend/middleware"
)

func ProviderRoutes(router *gin.Engine) {
	// Public routes
	router.GET("/api/providers", controllers.GetProviders)
	router.GET("/api/providers/", controllers.GetProviders) // Handle with trailing slash
	router.GET("/api/providers/:id", controllers.GetProvider)

	// Authorized routes
	authorized := router.Group("/api/providers")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/", controllers.CreateProvider)
		authorized.PUT("/:id", controllers.UpdateProvider)
		authorized.DELETE("/:id", controllers.DeleteProvider)
	}
}
