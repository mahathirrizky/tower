package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/user/tower-tracker-bima/backend/database"
	"github.com/user/tower-tracker-bima/backend/routes"
	"github.com/gin-gonic/contrib/static"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables directly")
	}

	// Initialize Database
	database.ConnectDatabase()
	database.SeedAdminUser()

	// Initialize Gin Router
	router := gin.Default()
	router.RedirectTrailingSlash = false // Disable automatic redirect for trailing slashes

	// CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins for debugging
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Serve static files for uploaded photos
	router.Static("/uploads", "./uploads")

	// Serve static frontend files from the 'dist' directory
	// This should be placed before API routes to ensure frontend assets are served first
	router.Use(static.Serve("/", static.LocalFile("../frontend/dist", true)))
	router.NoRoute(func(c *gin.Context) {
		c.File("../frontend/dist/index.html") // Serve index.html for all other routes (SPA fallback)
	})

	// Setup API routes
	log.Println("Registering Auth Routes...")
	routes.AuthRoutes(router)
	log.Println("Registering Provider Routes...")
	routes.ProviderRoutes(router)
	log.Println("Registering Tower Routes...")
	routes.TowerRoutes(router)
	log.Println("Registering Blankspot Routes...")
	routes.BlankspotRoutes(router)
	log.Println("All API routes registered.")

	// Start server
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}