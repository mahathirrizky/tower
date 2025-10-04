package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/user/tower-tracker-bima/backend/controllers"
	"github.com/user/tower-tracker-bima/backend/middleware"
)

func BlankspotRoutes(router *gin.Engine) {
	// Public routes (if any, though blankspots are likely admin-managed)
	router.GET("/api/blankspots", controllers.GetBlankspotAreas)
	router.GET("/api/blankspots/", controllers.GetBlankspotAreas) // Handle with trailing slash
	router.GET("/api/blankspots/:id", controllers.GetBlankspotArea)

	// Authorized routes
	authorized := router.Group("/api/blankspots")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/", controllers.CreateBlankspotArea)
		authorized.PUT("/:id", controllers.UpdateBlankspotArea)
		authorized.DELETE("/:id", controllers.DeleteBlankspotArea)
	}
}
