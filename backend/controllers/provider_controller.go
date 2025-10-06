package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/tower-tracker-bima/backend/database"
	"github.com/user/tower-tracker-bima/backend/helper"
	"github.com/user/tower-tracker-bima/backend/models"
)

type ProviderInput struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func CreateProvider(c *gin.Context) {
	log.Println("CreateProvider: Function entered")
	var input ProviderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("CreateProvider: ShouldBindJSON error: %v", err)
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("CreateProvider Input: %+v", input)

	provider := models.Provider{Name: input.Name, Address: input.Address}
	if err := database.DB.Create(&provider).Error; err != nil {
		log.Printf("CreateProvider: GORM Create error: %v", err)
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create provider: "+err.Error())
		return
	}
	log.Printf("CreateProvider: Provider successfully created in DB: %+v", provider)

	helper.SendSuccessResponse(c, http.StatusCreated, "Provider created successfully", provider)
}

func GetProviders(c *gin.Context) {
	var providers []models.Provider
	database.DB.Find(&providers)
	helper.SendSuccessResponse(c, http.StatusOK, "Providers fetched successfully", providers)
}

func GetProvider(c *gin.Context) {
	var provider models.Provider
	if err := database.DB.First(&provider, c.Param("id")).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Provider not found")
		return
	}
	helper.SendSuccessResponse(c, http.StatusOK, "Provider fetched successfully", provider)
}

func UpdateProvider(c *gin.Context) {
	var provider models.Provider
	if err := database.DB.First(&provider, c.Param("id")).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Provider not found")
		return
	}

	var input ProviderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updateData := map[string]interface{}{
		"Name":    input.Name,
		"Address": input.Address,
	}
	database.DB.Model(&provider).Updates(updateData)
	helper.SendSuccessResponse(c, http.StatusOK, "Provider updated successfully", provider)
}

func DeleteProvider(c *gin.Context) {
	var provider models.Provider
	// Use Unscoped to find the record even if it's soft-deleted
	if err := database.DB.Unscoped().First(&provider, c.Param("id")).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Provider not found")
		return
	}

	// Use Unscoped to permanently delete the record
	database.DB.Unscoped().Delete(&provider)
	helper.SendSuccessResponse(c, http.StatusOK, "Provider permanently deleted", nil)
}