package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/user/tower-tracker-bima/backend/controllers"
	"github.com/user/tower-tracker-bima/backend/middleware"
)

func TowerRoutes(router *gin.Engine) {
	// Public routes
	router.GET("/api/towers", controllers.GetTowers)
	router.GET("/api/towers/:id", controllers.GetTower)

	// Authorized routes
	authorized := router.Group("/api/towers")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("", controllers.CreateTower)
		authorized.PUT("/:id", controllers.UpdateTower)
		authorized.DELETE("/:id", controllers.DeleteTower)

		// New routes for timeline events
		authorized.PUT("/:id/ownership", controllers.ChangeOwnership)
		authorized.PUT("/:id/relocate", controllers.RelocateTower)
		authorized.PUT("/:id/dismantle", controllers.DismantleTower)
		authorized.GET("/:id/history", controllers.GetTowerHistory)
	}
}
