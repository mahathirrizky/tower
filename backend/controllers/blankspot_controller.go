package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/tower-tracker-bima/backend/database"
	"github.com/user/tower-tracker-bima/backend/helper"
	"github.com/user/tower-tracker-bima/backend/models"
)

// BlankspotAreaInput defines the structure for creating/updating a blankspot area
type BlankspotAreaInput struct {
	Name        string `json:"name" binding:"required"`
	Kelurahan   string `json:"kelurahan" binding:"required"`
	Coordinates string `json:"coordinates" binding:"required"` // JSON string of [[lat, lon], ...]
	Type        string `json:"type" binding:"required"`        // e.g., "Blankspot", "Weak Signal"
	Color       string `json:"color" binding:"required"`       // e.g., "#FF0000", "#FFFF00"
}

// CreateBlankspotArea handles the creation of a new blankspot area
func CreateBlankspotArea(c *gin.Context) {
	var input BlankspotAreaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	blankspotArea := models.BlankspotArea{
		Name:        input.Name,
		Kelurahan:   input.Kelurahan,
		Coordinates: input.Coordinates,
		Type:        input.Type,
		Color:       input.Color,
	}

	if err := database.DB.Create(&blankspotArea).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create blankspot area: "+err.Error())
		return
	}

	helper.SendSuccessResponse(c, http.StatusCreated, "Blankspot area created successfully", blankspotArea)
}

// GetBlankspotAreas handles fetching all blankspot areas
func GetBlankspotAreas(c *gin.Context) {
	var blankspotAreas []models.BlankspotArea
	database.DB.Find(&blankspotAreas)
	helper.SendSuccessResponse(c, http.StatusOK, "Blankspot areas fetched successfully", blankspotAreas)
}

// GetBlankspotArea handles fetching a single blankspot area by ID
func GetBlankspotArea(c *gin.Context) {
	var blankspotArea models.BlankspotArea
	if err := database.DB.First(&blankspotArea, c.Param("id")).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Blankspot area not found")
		return
	}
	helper.SendSuccessResponse(c, http.StatusOK, "Blankspot area fetched successfully", blankspotArea)
}

// UpdateBlankspotArea handles updating an existing blankspot area
func UpdateBlankspotArea(c *gin.Context) {
	var blankspotArea models.BlankspotArea
	if err := database.DB.First(&blankspotArea, c.Param("id")).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Blankspot area not found")
		return
	}

	var input BlankspotAreaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updateData := map[string]interface{}{
		"Name":        input.Name,
		"Kelurahan":   input.Kelurahan,
		"Coordinates": input.Coordinates,
		"Type":        input.Type,
		"Color":       input.Color,
	}
	database.DB.Model(&blankspotArea).Updates(updateData)
	helper.SendSuccessResponse(c, http.StatusOK, "Blankspot area updated successfully", blankspotArea)
}

// DeleteBlankspotArea handles deleting a blankspot area
func DeleteBlankspotArea(c *gin.Context) {
	var blankspotArea models.BlankspotArea
	if err := database.DB.Unscoped().First(&blankspotArea, c.Param("id")).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "Blankspot area not found")
		return
	}

	database.DB.Unscoped().Delete(&blankspotArea)
	helper.SendSuccessResponse(c, http.StatusOK, "Blankspot area permanently deleted", nil)
}
