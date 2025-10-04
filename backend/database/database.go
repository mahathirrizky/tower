package database

import (
	"log"

	"github.com/user/tower-tracker-bima/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("tower_tracker.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	err = database.AutoMigrate(&models.User{}, &models.Provider{}, &models.Tower{}, &models.BlankspotArea{}, &models.TowerEvent{})
	if err != nil {
		log.Fatal("Failed to migrate database!", err)
	}

	DB = database
	log.Println("Database connection and migration successful.")
}

func SeedAdminUser() {
	var user models.User
	DB.Where("email = ?", "admin@example.com").First(&user)

	if user.ID == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		adminUser := models.User{
			Email:    "admin@example.com",
			Password: string(hashedPassword),
		}
		DB.Create(&adminUser)
		log.Println("Admin user 'admin@example.com' seeded with password 'password'")
	} else {
		log.Println("Admin user already exists.")
	}
}
