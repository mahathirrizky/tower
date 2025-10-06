package controllers

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chai2010/webp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/tower-tracker-bima/backend/database"
	"github.com/user/tower-tracker-bima/backend/helper"
	"github.com/user/tower-tracker-bima/backend/models"
)

// processAndSaveWebP handles decoding an uploaded image, converting it to WebP, and saving it.
func processAndSaveWebP(file *multipart.FileHeader) (string, error) {
	// 1. Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// 2. Decode the image
	img, _, err := image.Decode(src)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// 3. Generate a unique filename with .webp extension
	newFileName := uuid.New().String() + ".webp"
	dstPath := filepath.Join("uploads", newFileName)

	// 4. Create the destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// 5. Encode the image to WebP and save
	if err := webp.Encode(dstFile, img, &webp.Options{Quality: 90}); err != nil {
		return "", fmt.Errorf("failed to encode image to webp: %w", err)
	}

	// 6. Return the URL path
	return fmt.Sprintf("/uploads/%s", newFileName), nil
}

func CreateTower(c *gin.Context) {
	// Parse form data
	latStr := c.PostForm("latitude")
	lonStr := c.PostForm("longitude")
	kelurahan := c.PostForm("kelurahan")
	kecamatan := c.PostForm("kecamatan")
	address := c.PostForm("address")
	tipe := c.PostForm("tipe")
	tinggiStr := c.PostForm("tinggi")
	providerIDs := c.PostFormArray("provider_ids")

	log.Printf("CreateTower: Received form data: lat=%s, lon=%s, kelurahan=%s, kecamatan=%s, tipe=%s, tinggi=%s, providers=%v", latStr, lonStr, kelurahan, kecamatan, tipe, tinggiStr, providerIDs)

	c.Request.ParseMultipartForm(10 << 20) // 10 MB
	log.Printf("Full form data: %v", c.Request.Form)

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid latitude format")
		return
	}

	longitude, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid longitude format")
		return
	}

	var tinggi float64
	if tinggiStr != "" {
		tinggi, err = strconv.ParseFloat(tinggiStr, 64)
		if err != nil {
			helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid tinggi format")
			return
		}
	}

	// Handle file upload and convert to WebP
	var photoURL string
	file, err := c.FormFile("photo")
	if err == nil { // Photo is provided
		log.Println("CreateTower: Photo file found, processing...")
		photoURL, err = processAndSaveWebP(file)
		if err != nil {
			helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	} else if err != http.ErrMissingFile {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Error processing photo upload")
		return
	}

	// Find providers
	var providers []*models.Provider
	if len(providerIDs) > 0 {
		var providerModels []models.Provider
		if err := database.DB.Where(providerIDs).Find(&providerModels).Error; err != nil {
			helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to find providers")
			return
		}
		for i := range providerModels {
			providers = append(providers, &providerModels[i])
		}
	}

	tower := models.Tower{
		Latitude:  latitude,
		Longitude: longitude,
		Kelurahan: kelurahan,
		Kecamatan: kecamatan,
		Address:   address,
		Tinggi:    tinggi,
		Tipe:      tipe,
		PhotoURL:  photoURL, // Use the new WebP photo URL
		Providers: providers,
		Status:    "active", // Default status
	}

	log.Printf("CreateTower: Attempting to create tower object: %+v", tower)

	dbResult := database.DB.Create(&tower)
	if dbResult.Error != nil {
		log.Printf("CreateTower: DB.Create FAILED: %v", dbResult.Error)
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create tower: "+dbResult.Error.Error())
		return
	}

	log.Printf("CreateTower: DB.Create SUCCEEDED. Rows affected: %d. New Tower ID: %d", dbResult.RowsAffected, tower.ID)

	// Create TowerEvent for creation
	if err := createTowerEvent(c, tower.ID, "Created", "Tower initially created.", nil, gin.H{
		"latitude":  tower.Latitude,
		"longitude": tower.Longitude,
		"kelurahan": tower.Kelurahan,
		"kecamatan": tower.Kecamatan,
		"tinggi":    tower.Tinggi,
		"tipe":      tower.Tipe,
		"photo_url": tower.PhotoURL,
		"providers": providers,
		"status":    tower.Status,
	}); err != nil {
		fmt.Printf("Failed to create tower event for creation: %v\n", err)
	}

	helper.SendSuccessResponse(c, http.StatusCreated, "Tower created successfully", tower)
}

func GetTowers(c *gin.Context) {
	var towers []models.Tower
	database.DB.Preload("Providers").Find(&towers)
	helper.SendSuccessResponse(c, http.StatusOK, "Towers fetched successfully", towers)
}

func GetTower(c *gin.Context) {
	var tower models.Tower
	if err := database.DB.Preload("Providers").First(&tower, c.Param("id")).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Tower not found")
		return
	}
	helper.SendSuccessResponse(c, http.StatusOK, "Tower fetched successfully", tower)
}

func UpdateTower(c *gin.Context) {
	var tower models.Tower
	if err := database.DB.First(&tower, c.Param("id")).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Tower not found")
		return
	}

	// Store old data for event logging
	oldData := gin.H{
		"kelurahan": tower.Kelurahan,
		"kecamatan": tower.Kecamatan,
		"tinggi":    tower.Tinggi,
		"tipe":      tower.Tipe,
		"photo_url": tower.PhotoURL,
	}

	// Parse form data for details
	kelurahan := c.PostForm("kelurahan")
	kecamatan := c.PostForm("kecamatan")
	address := c.PostForm("address")
	tipe := c.PostForm("tipe")
	tinggiStr := c.PostForm("tinggi")

	var tinggi float64
	if tinggiStr != "" {
		var err error
		tinggi, err = strconv.ParseFloat(tinggiStr, 64)
		if err != nil {
			helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid tinggi format")
			return
		}
	}

	// Handle file upload (optional)
	file, err := c.FormFile("photo")
	if err == nil { // Photo provided, update it
		if tower.PhotoURL != "" {
			oldPhotoPath := filepath.Join(".", tower.PhotoURL)
			if err := os.Remove(oldPhotoPath); err != nil {
				fmt.Printf("Failed to delete old photo file %s: %v\n", oldPhotoPath, err)
			}
		}
		newPhotoURL, err := processAndSaveWebP(file)
		if err != nil {
			helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		tower.PhotoURL = newPhotoURL
	} else if err != http.ErrMissingFile {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Error processing photo upload")
		return
	}

	// Update tower detail fields
	tower.Kelurahan = kelurahan
	tower.Kecamatan = kecamatan
	tower.Address = address
	tower.Tipe = tipe
	if tinggiStr != "" { // Only update if provided
		tower.Tinggi = tinggi
	}

	if err := database.DB.Save(&tower).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update tower: "+err.Error())
		return
	}

	// Create TowerEvent for details update
	newData := gin.H{
		"kelurahan": tower.Kelurahan,
		"kecamatan": tower.Kecamatan,
		"tinggi":    tower.Tinggi,
		"tipe":      tower.Tipe,
		"photo_url": tower.PhotoURL,
	}

	if fmt.Sprintf("%v", oldData) != fmt.Sprintf("%v", newData) {
		if err := createTowerEvent(c, tower.ID, "DetailsUpdate", "Tower details were updated.", oldData, newData); err != nil {
			fmt.Printf("Failed to create tower event for details update: %v\n", err)
		}
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Tower updated successfully", tower)
}

func DeleteTower(c *gin.Context) {
	var tower models.Tower
	if err := database.DB.First(&tower, c.Param("id")).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Tower not found")
		return
	}

	// Store old data for event logging
	oldStatus := tower.Status

	// Delete the associated photo file
	if tower.PhotoURL != "" {
		photoPath := filepath.Join(".", tower.PhotoURL)
		if err := os.Remove(photoPath); err != nil {
			fmt.Printf("Failed to delete photo file %s: %v\n", photoPath, err)
		}
	}

	database.DB.Delete(&tower)

	// Create TowerEvent for dismantling
	if err := createTowerEvent(c, tower.ID, "Dismantled", fmt.Sprintf("Tower deleted (status changed from %s to dismantled)", oldStatus), gin.H{"status": oldStatus}, gin.H{"status": "dismantled"}); err != nil {
		fmt.Printf("Failed to create tower event for dismantling: %v\n", err)
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Tower deleted successfully", nil)
}

// Helper to compare provider slices
func compareProviders(oldProviders, newProviders []*models.Provider) bool {
	if len(oldProviders) != len(newProviders) {
		return false
	}
	oldMap := make(map[uint]bool)
	for _, p := range oldProviders {
		oldMap[p.ID] = true
	}
	for _, p := range newProviders {
		if !oldMap[p.ID] {
			return false
		}
	}
	return true
}

// Helper to get provider names
func getProviderNames(providers []*models.Provider) []string {
	names := []string{}
	for _, p := range providers {
		names = append(names, p.Name)
	}
	return names
}

// Helper to create TowerEvent
func createTowerEvent(c *gin.Context, towerID uint, eventType, description string, oldData, newData interface{}) error {
	userID := c.MustGet("user_id").(uint)

	oldDataJSON, _ := json.Marshal(oldData)
	newDataJSON, _ := json.Marshal(newData)

	towerEvent := models.TowerEvent{
		TowerID:     towerID,
		EventType:   eventType,
		Timestamp:   time.Now(),
		Description: description,
		OldData:     string(oldDataJSON),
		NewData:     string(newDataJSON),
		UserID:      userID,
	}
	return database.DB.Create(&towerEvent).Error
}

// ChangeOwnershipInput defines input for changing tower ownership
type ChangeOwnershipInput struct {
	NewProviderID uint `json:"new_provider_id" binding:"required"`
}

// ChangeOwnership handles changing the provider of a tower
func ChangeOwnership(c *gin.Context) {
	towerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid tower ID")
		return
	}

	var input ChangeOwnershipInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var tower models.Tower
	if err := database.DB.Preload("Providers").First(&tower, towerID).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Tower not found")
		return
	}

	var newProvider models.Provider
	if err := database.DB.First(&newProvider, input.NewProviderID).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "New provider not found")
		return
	}

	// Store old provider info for event
	oldProviders := tower.Providers
	oldProviderNames := []string{}
	for _, p := range oldProviders {
		oldProviderNames = append(oldProviderNames, p.Name)
	}

	// Update providers association
	if err := database.DB.Model(&tower).Association("Providers").Replace(&newProvider); err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to change tower ownership")
		return
	}

	// Create TowerEvent
	description := fmt.Sprintf("Ownership changed from %s to %s", oldProviderNames, newProvider.Name)
	oldData := gin.H{"providers": oldProviderNames}
	newData := gin.H{"providers": []string{newProvider.Name}}
	if err := createTowerEvent(c, tower.ID, "OwnershipChange", description, oldData, newData); err != nil {
		fmt.Printf("Failed to create tower event: %v\n", err)
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Tower ownership changed successfully", tower)
}

// RelocateTowerInput defines input for relocating a tower
type RelocateTowerInput struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// RelocateTower handles changing the location of a tower
func RelocateTower(c *gin.Context) {
	towerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid tower ID")
		return
	}

	var input RelocateTowerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var tower models.Tower
	if err := database.DB.First(&tower, towerID).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Tower not found")
		return
	}

	// Store old location for event
	oldLatitude := tower.Latitude
	oldLongitude := tower.Longitude

	// Update tower location
	tower.Latitude = input.Latitude
	tower.Longitude = input.Longitude

	if err := database.DB.Save(&tower).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to relocate tower")
		return
	}

	// Create TowerEvent
	description := fmt.Sprintf("Tower relocated from (%.6f, %.6f) to (%.6f, %.6f)", oldLatitude, oldLongitude, tower.Latitude, tower.Longitude)
	oldData := gin.H{"latitude": oldLatitude, "longitude": oldLongitude}
	newData := gin.H{"latitude": tower.Latitude, "longitude": tower.Longitude}
	if err := createTowerEvent(c, tower.ID, "Relocation", description, oldData, newData); err != nil {
		fmt.Printf("Failed to create tower event: %v\n", err)
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Tower relocated successfully", tower)
}

// DismantleTower handles marking a tower as dismantled
func DismantleTower(c *gin.Context) {
	towerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid tower ID")
		return
	}

	var tower models.Tower
	if err := database.DB.First(&tower, towerID).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Tower not found")
		return
	}

	// Store old status for event
	oldStatus := tower.Status

	// Update tower status
	tower.Status = "dismantled"

	if err := database.DB.Save(&tower).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to dismantle tower")
		return
	}

	// Create TowerEvent
	description := fmt.Sprintf("Tower status changed from %s to %s", oldStatus, tower.Status)
	oldData := gin.H{"status": oldStatus}
	newData := gin.H{"status": tower.Status}
	if err := createTowerEvent(c, tower.ID, "Dismantled", description, oldData, newData); err != nil {
		fmt.Printf("Failed to create tower event: %v\n", err)
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Tower dismantled successfully", tower)
}

// GetTowerHistory retrieves the event history for a specific tower
func GetTowerHistory(c *gin.Context) {
	towerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid tower ID")
		return
	}

	var towerEvents []models.TowerEvent
	if err := database.DB.Where("tower_id = ?", towerID).Order("timestamp desc").Find(&towerEvents).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch tower history")
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Tower history fetched successfully", towerEvents)
}
